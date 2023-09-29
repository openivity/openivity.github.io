package main

import (
	"fmt"
	"syscall/js"

	"github.com/muktihari/openactivity-fit/activity"
)

func main() {
	fmt.Println("Fit Service WebAssembly Instantiated")
	js.Global().Set("decode", js.FuncOf(activity.Decode))
	select {} // never exit
}
