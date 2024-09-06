// Copyright (C) 2023 Openivity

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"github.com/muktihari/fit/profile/typedef"
	"github.com/openivity/activity-service/activity"
	"github.com/openivity/activity-service/activity/fit"
	"github.com/openivity/activity-service/activity/gpx"
	"github.com/openivity/activity-service/activity/tcx"
	"github.com/openivity/activity-service/mem"
	"github.com/openivity/activity-service/service"
	"github.com/openivity/activity-service/service/spec"
	"github.com/openivity/activity-service/strutils"
	"golang.org/x/exp/slices"
)

// decodedActivities is previously decoded activities.
//
// NOTE: We cache this for faster encoding process, serializing and deserializing data
// in the current Go WebAssembly implementation is expensive, as the [syscall/js] library
// is still considered EXPERIMENTAL. This is safe since every WebAssembly Instance is isolated.
var decodedActivities = []activity.Activity{}

//go:embed manufacturers.json
var manufacturerJson []byte

func main() {
	preproc := activity.NewPreprocessor()

	svc := service.New(
		fit.NewDecodeEncoder(preproc),
		gpx.NewDecodeEncoder(preproc),
		tcx.NewDecodeEncoder(preproc),
		makeManufacturers(),
	)

	js.Global().Set("decode", createDecodeFunc(svc))
	js.Global().Set("encode", createEncodeFunc(svc))
	js.Global().Set("manufacturerList", createManufacturerListFunc(svc))
	js.Global().Set("sportList", createSportListFunc(svc))

	// Add shutdown hook
	quitc := make(chan struct{})
	js.Global().Set("shutdown", js.FuncOf(func(this js.Value, args []js.Value) any {
		close(quitc)
		return nil
	}))

	fmt.Println("WebAssembly: Activity Service Instantiated")
	<-quitc
	fmt.Println("WebAssembly: Activity Service Exited!")
}

func makeManufacturers() map[typedef.Manufacturer]activity.Manufacturer {
	manufacturers := make(map[typedef.Manufacturer]activity.Manufacturer)

	var manufacturerJsonData map[string]activity.Manufacturer
	if err := json.Unmarshal(manufacturerJson, &manufacturerJsonData); err != nil {
		// Only happen if manufacturers.json is corrupted on build, less likely to happen. Let's just log it.
		fmt.Printf("could not make manufactures mapping: %v\n", err)
		return manufacturers
	}

	manufacturerList := typedef.ListManufacturer()

	for i := range manufacturerList {
		manufacturer := activity.Manufacturer{
			ID:   manufacturerList[i],
			Name: strutils.ToTitle(manufacturerList[i].String()),
		}

		if manufacturer.ID == typedef.ManufacturerGarmin {
			garminProductIDs := typedef.ListGarminProduct()
			for j := range garminProductIDs {
				product := activity.Product{
					ID:   garminProductIDs[j].Uint16(),
					Name: strutils.ToTitle(garminProductIDs[j].String()),
				}
				manufacturer.Products = append(manufacturer.Products, product)
			}
		}

		if m, ok := manufacturerJsonData[strconv.FormatUint(uint64(manufacturer.ID), 10)]; ok {
			manufacturer.Name = m.Name
			manufacturer.Products = append(manufacturer.Products, m.Products...)
		}

		slices.SortFunc(manufacturer.Products, func(a, b activity.Product) int {
			if strings.ToLower(a.Name) < strings.ToLower(b.Name) {
				return -1
			}
			return 1
		})

		manufacturers[manufacturer.ID] = manufacturer
	}

	return manufacturers
}

func createDecodeFunc(s *service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := args[0] // input is an Array<Uint8Array>
		if input.Length() == 0 {
			return "{\"err\":\"no input is passed.\"}"
		}

		rs := make([]io.Reader, input.Length())

		for i := 0; i < input.Length(); i++ {
			b := make([]byte, input.Index(i).Length())
			js.CopyBytesToGo(b, input.Index(i))
			rs[i] = bytes.NewReader(b)
		}

		result := s.Decode(context.Background(), rs)
		decodedActivities = result.Activities

		buf := mem.GetBuffer()
		defer mem.PutBuffer(buf)

		b := result.MarshalAppendJSON(buf.Bytes())

		return string(b)
	})
}

func createEncodeFunc(svc *service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := args[0] // input is an JSON string
		if input.Length() == 0 {
			return "{\"err\":\"no input is passed.\"}"
		}

		begin := time.Now()
		b := make([]byte, input.Length())
		js.CopyBytesToGo(b, input)

		var encodeSpec spec.Encode
		if err := json.Unmarshal(b, &encodeSpec); err != nil {
			return "{\"err\":\"could not unmarshal input\"}"
		}

		encodeSpec.Activities = cloneActivities(decodedActivities)
		elapsed := time.Since(begin)

		result := svc.Encode(context.Background(), encodeSpec)
		result.DeserializeInputTook = elapsed

		buf := mem.GetBuffer()
		defer mem.PutBuffer(buf)

		b = result.MarshalAppendJSON(buf.Bytes())

		return string(b)
	})
}

func createManufacturerListFunc(svc *service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		manufacturerList := svc.ManufacturerList()
		b := manufacturerList.MarshalAppendJSON(make([]byte, 0, 50<<10))
		return string(b)
	})
}

func createSportListFunc(svc *service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		sportList := svc.SportList()
		b := sportList.MarshalAppendJSON(make([]byte, 0, 8<<10))
		return string(b)
	})
}

// cloneActivities clones activities so each encode invocation has isolated activities data.
func cloneActivities(activities []activity.Activity) []activity.Activity {
	activities = slices.Clone(activities)
	for i := range activities {
		sessions := slices.Clone(activities[i].Sessions)
		for j := range sessions {
			base := *sessions[j].Session
			sessions[j].Session = &base

			records := slices.Clone(sessions[j].Records)
			for k := range records {
				base := *records[k].Record
				records[k].Record = &base
			}
			sessions[j].Records = records

			laps := slices.Clone(sessions[j].Laps)
			for k := range laps {
				base := *laps[k].Lap
				laps[k].Lap = &base
			}
			sessions[j].Laps = laps
		}
		activities[i].Sessions = sessions
	}

	return activities
}
