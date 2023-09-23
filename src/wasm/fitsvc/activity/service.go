package activity

import (
	"bytes"
	"context"
	"fmt"
	"syscall/js"
	"time"

	"github.com/muktihari/fit/decoder"
)

type Result struct {
	Err          string        `json:"err"`
	Feature      *Feature      `json:"feature"`
	ActivityFile *ActivityFile `json:"activityFile"`
}

func (m Result) ToMap() map[string]any {
	return map[string]any{
		"err":          m.Err,
		"feature":      m.Feature.ToMap(),
		"activityFile": m.ActivityFile.ToMap(),
	}
}

func Decode() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		data := args[0]
		dest := make([]byte, data.Length())
		js.CopyBytesToGo(dest, data)

		lis := NewListener()
		dec := decoder.New(bytes.NewReader(dest),
			decoder.WithMesgListener(lis),
			decoder.WithBroadcastOnly(),
		)

		begin := time.Now()
		_, err := dec.Decode(context.Background())
		if err != nil {
			return Result{Err: fmt.Sprintf("could not decode: %s", err)}.ToMap()
		}
		lis.Wait()

		fmt.Printf("fitsvc: decode took: %s\n", time.Since(begin))
		begin = time.Now()
		res := Result{
			Feature:      lis.Feature(),
			ActivityFile: lis.ActivityFile(),
		}.ToMap()

		fmt.Printf("fitsvc: serialization took: %s\n", time.Since(begin))

		return res

	})
}
