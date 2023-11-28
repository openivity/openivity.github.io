package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"syscall/js"

	"github.com/muktihari/fit/profile/typedef"
	"github.com/muktihari/openactivity-fit/activity/fit"
	"github.com/muktihari/openactivity-fit/activity/gpx"
	"github.com/muktihari/openactivity-fit/activity/tcx"
	"github.com/muktihari/openactivity-fit/kit"
	"github.com/muktihari/openactivity-fit/preprocessor"
	"github.com/muktihari/openactivity-fit/service"
	"github.com/muktihari/openactivity-fit/service/result"
	"golang.org/x/exp/slices"
)

//go:embed manufacturers.json
var manufacturerJson []byte

func main() {
	manufacturers, err := makeManufacturers()
	if err != nil {
		fmt.Printf("could not make manufactures mapping: %v\n", err)
	}

	preproc := preprocessor.New()

	fs := fit.NewService(preproc, manufacturers)
	gs := gpx.NewService(preproc)
	ts := tcx.NewService(preproc)

	s := service.New(fs, gs, ts, manufacturers)

	js.Global().Set("decode", createDecodeFunc(s))
	js.Global().Set("encode", createEncodeFunc(s))
	js.Global().Set("manufacturerList", createManufacturerListFunc(s))

	fmt.Println("WebAssembly: Activity Service Instantiated")
	select {} // never exit
}

func makeManufacturers() (manufacturers map[uint16]fit.Manufacturer, err error) {
	manufacturers = make(map[uint16]fit.Manufacturer)

	var source map[string]fit.Manufacturer
	if err = json.Unmarshal(manufacturerJson, &source); err != nil {
		return
	}

	manufacturerIDs := typedef.ListManufacturer()
	garminProductIDs := typedef.ListGarminProduct()

	for i := range manufacturerIDs {
		manufacturerID := manufacturerIDs[i]
		manufacturer := fit.Manufacturer{
			ID:   uint16(manufacturerID),
			Name: kit.FormatTitle(manufacturerID.String()),
		}

		if manufacturer.ID == uint16(typedef.ManufacturerGarmin) {
			for j := range garminProductIDs {
				product := fit.ManufacturerProduct{
					ID:   uint16(garminProductIDs[j]),
					Name: kit.FormatTitle(garminProductIDs[j].String()),
				}
				manufacturer.Products = append(manufacturer.Products, product)
			}
		}

		if m, ok := source[strconv.FormatUint(uint64(manufacturer.ID), 10)]; ok {
			manufacturer.Name = m.Name
			manufacturer.Products = m.Products
		}

		slices.SortFunc(manufacturer.Products, func(a, b fit.ManufacturerProduct) int {
			if a.Name < b.Name {
				return -1
			}
			return 1
		})

		manufacturers[manufacturer.ID] = manufacturer
	}

	return
}

func createDecodeFunc(s service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := args[0] // input is an Array<Uint8Array>
		if input.Length() == 0 {
			return result.Decode{Err: fmt.Errorf("no input is passed")}.ToMap()
		}

		rs := make([]io.Reader, input.Length())

		for i := 0; i < input.Length(); i++ {
			b := make([]byte, input.Index(i).Length())
			js.CopyBytesToGo(b, input.Index(i))
			rs[i] = bytes.NewReader(b)
		}

		result := s.Decode(context.Background(), rs)

		return result.ToMap()
	})
}

func createEncodeFunc(s service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		result := s.Encode(context.Background(), nil)
		return result.ToMap()
	})
}

func createManufacturerListFunc(s service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		return s.ManufacturerList().ToMap()
	})
}
