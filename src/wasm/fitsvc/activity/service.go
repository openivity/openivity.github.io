package activity

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"syscall/js"
	"time"

	"github.com/muktihari/fit/decoder"
)

type Result struct {
	Err          string       `json:"err"`
	Feature      Feature      `json:"feature"`
	ActivityFile ActivityFile `json:"activityFile"`
}

func Decode() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		begin := time.Now()

		data := args[0]
		dest := make([]byte, data.Length())
		js.CopyBytesToGo(dest, data)

		lis := NewListener()
		dec := decoder.New(bufio.NewReader(bytes.NewReader(dest)),
			decoder.WithMesgListener(lis),
			decoder.WithBroadcastOnly(),
		)

		_, err := dec.Decode(context.Background())
		if err != nil {
			return encodeResult(Result{Err: fmt.Sprintf("could not decode: %s", err)})
		}
		lis.Wait()

		fmt.Printf("decode took: %s\n", time.Since(begin))

		return encodeResult(Result{
			Feature:      lis.Feature(),
			ActivityFile: lis.ActivityFile(),
		})

	})
}

func encodeResult(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}
