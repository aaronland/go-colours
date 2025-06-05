//go:build wasmjs
package main

import (
	"syscall/js"
	"log"
	
	"github.com/aaronland/go-colours/wasm"
)

func main() {

	extrude_func := wasm.ExtrudeFunc()
	defer extrude_func.Release()

	js.Global().Set("colours_extrude", extrude_func)

	c := make(chan struct{}, 0)

	log.Println("WASM colours functions initialized")
	<-c
}
