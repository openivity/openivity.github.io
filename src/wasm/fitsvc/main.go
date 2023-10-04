package main

import (
	"fmt"
	"syscall/js"

	"github.com/muktihari/openactivity-fit/activity"
)

func main() {
	js.Global().Set("decode", activity.Decode())

	fmt.Println("Fit Service WebAssembly Instantiated")
	select {} // never exit
}
