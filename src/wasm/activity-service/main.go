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
	"sync"
	"syscall/js"
	"time"

	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/openactivity-fit/activity"
	"github.com/muktihari/openactivity-fit/activity/fit"
	"github.com/muktihari/openactivity-fit/activity/gpx"
	"github.com/muktihari/openactivity-fit/activity/tcx"
	"github.com/muktihari/openactivity-fit/kit"
	"github.com/muktihari/openactivity-fit/preprocessor"
	"github.com/muktihari/openactivity-fit/service"
	"github.com/muktihari/openactivity-fit/service/spec"
	"golang.org/x/exp/slices"
)

// Since we don't require concurrency on functions invocation, let's cache previously decoded activities
// for a faster encoding process. Serializing and deserializing data is quite expensive in the current
// WebAssembly specification, especially in Golang, as the [syscall/js] library is still considered EXPERIMENTAL.
//
// This implementation is similiar on how Linier Memory work in WebAssembly.
// ref:
//   - https://developer.mozilla.org/en-US/docs/WebAssembly/JavaScript_interface/Memory
//   - https://wasmbyexample.dev/examples/webassembly-linear-memory/webassembly-linear-memory.go.en-us.html
var decodedActivitiesCache = [1][]activity.Activity{} // 1 single process is allowed at a time.
const pointer = 0                                     // now, it's always pointing to zero
var mu sync.Mutex                                     // cache access lock

//go:embed manufacturers.json
var manufacturerJson []byte

func main() {
	manufacturers := makeManufacturers()

	preproc := preprocessor.New()

	fs := fit.NewService(preproc, manufacturers)
	gs := gpx.NewService(preproc)
	ts := tcx.NewService(preproc)

	s := service.New(fs, gs, ts, manufacturers)

	js.Global().Set("decode", createDecodeFunc(s))
	js.Global().Set("encode", createEncodeFunc(s))
	js.Global().Set("manufacturerList", createManufacturerListFunc(s))
	js.Global().Set("sportList", createSportListFunc(s))

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

// cache caches decoded result to be used on encoding later.
func cache(activities []activity.Activity) {
	mu.Lock()
	decodedActivitiesCache[pointer] = activities
	mu.Unlock()
}

// retrieveCache retrieves a cloned version of latest decoded value.
func retrieveCache(pointer byte) []activity.Activity {
	mu.Lock()
	cachedActivities := decodedActivitiesCache[pointer]
	activities := make([]activity.Activity, len(cachedActivities))
	for i := range cachedActivities {
		activities[i] = *cachedActivities[i].Clone()
	}
	mu.Unlock()
	return activities
}

func makeManufacturers() map[uint16]fit.Manufacturer {
	manufacturers := make(map[uint16]fit.Manufacturer)

	var source map[string]fit.Manufacturer
	if err := json.Unmarshal(manufacturerJson, &source); err != nil {
		// Only happen if manufacturers.json is corrupted on build, less likely to happen. Let's just log it.
		fmt.Printf("Could not make manufactures mapping: %v\n", err)
		return manufacturers
	}

	manufacturerIDs := typedef.ListManufacturer()
	garminProductIDs := typedef.ListGarminProduct()

	for i := range manufacturerIDs {
		if manufacturerIDs[i] == typedef.ManufacturerInvalid {
			continue
		}

		manufacturerID := manufacturerIDs[i]
		manufacturer := fit.Manufacturer{
			ID:   uint16(manufacturerID),
			Name: kit.FormatTitle(manufacturerID.String()),
		}

		if manufacturer.ID == uint16(typedef.ManufacturerGarmin) {
			for j := range garminProductIDs {
				if garminProductIDs[j] == typedef.GarminProductInvalid {
					continue
				}

				product := fit.ManufacturerProduct{
					ID:   uint16(garminProductIDs[j]),
					Name: kit.FormatTitle(garminProductIDs[j].String()),
				}
				manufacturer.Products = append(manufacturer.Products, product)
			}
		}

		if m, ok := source[strconv.FormatUint(uint64(manufacturer.ID), 10)]; ok {
			manufacturer.Name = m.Name
			if len(manufacturer.Products) == 0 {
				manufacturer.Products = m.Products
			} else {
				for i := range m.Products {
					mp := m.Products[i]
					var exist bool
					for j := range manufacturer.Products {
						p := &manufacturer.Products[j]
						if p.ID == mp.ID {
							p.Name = mp.Name // Rename product name
							exist = true
							break
						}
					}
					if !exist {
						manufacturer.Products = append(manufacturer.Products, mp)
					}
				}
			}
		}

		slices.SortFunc(manufacturer.Products, func(a, b fit.ManufacturerProduct) int {
			if strings.ToLower(a.Name) < strings.ToLower(b.Name) {
				return -1
			}
			return 1
		})

		manufacturers[manufacturer.ID] = manufacturer
	}

	return manufacturers
}

func createDecodeFunc(s service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := args[0] // input is an Array<Uint8Array>
		if input.Length() == 0 {
			return map[string]any{"err": fmt.Errorf("no input is passed")}
		}

		rs := make([]io.Reader, input.Length())

		for i := 0; i < input.Length(); i++ {
			b := make([]byte, input.Index(i).Length())
			js.CopyBytesToGo(b, input.Index(i))
			rs[i] = bytes.NewReader(b)
		}

		result := s.Decode(context.Background(), rs)

		cache(result.Activities)

		b, _ := result.MarshalJSON()

		return string(b)
	})
}

func createEncodeFunc(s service.Service) js.Func {
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

		encodeSpec.Activities = retrieveCache(pointer)

		elapsed := time.Since(begin)

		result := s.Encode(context.Background(), encodeSpec)
		result.DeserializeInputTook = elapsed

		b, _ = result.MarshalJSON()

		return string(b)
	})
}

func createManufacturerListFunc(s service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		manufacturerList := s.ManufacturerList()
		b, _ := manufacturerList.MarshalJSON()
		return string(b)
	})
}

func createSportListFunc(s service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		sportList := s.SportList()
		b, _ := sportList.MarshalJSON()
		return string(b)
	})
}
