package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"syscall/js"

	"github.com/muktihari/openactivity-fit/activity/fit"
	"github.com/muktihari/openactivity-fit/activity/gpx"
	"github.com/muktihari/openactivity-fit/activity/tcx"
	"github.com/muktihari/openactivity-fit/preprocessor"
	"github.com/muktihari/openactivity-fit/service"
)

func main() {
	preproc := preprocessor.New()

	fs := fit.NewService(preproc)
	gs := gpx.NewService(preproc)
	ts := tcx.NewService(preproc)

	s := service.New(fs, gs, ts)

	js.Global().Set("decode", Decode(s))

	fmt.Println("WebAssembly: Activity Service Instantiated")
	select {} // never exit
}

func Decode(s service.Service) js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		input := args[0] // input is an Array<Uint8Array>
		if input.Length() == 0 {
			return service.Result{Err: fmt.Errorf("no input is passed")}.ToMap()
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
