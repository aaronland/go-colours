//go:build wasmjs
package wasm

import (
	"strings"
	"fmt"
	"bytes"
	"bufio"
	"syscall/js"
	"encoding/json"
	"encoding/base64"
	
	"github.com/aaronland/go-colours/extrude"
)

func ExtrudeFunc () js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		b64_im := args[0].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			im_data, err := base64.StdEncoding.DecodeString(im_b64)
			
			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to decode image body, %w", err))
				return nil
			}

			im_r := bytes.NewReader(im_data)
			
			im, _, err := image.Decode(im_r)

			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to decode image, %w", err))
				return nil
			}
			
			rsp, err := extrude.ExtrudeImages(ctx, extrude_opts, im)

			im, _, err := image.Decode(im_r)

			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to extrude image colours, %w", err))
				return nil
			}
			
			enc_rsp, err := json.Marshal(rsp)

			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to marshal image colours, %w", err))
				return nil
			}
			
			resolve.Invoke(string(enc_rsp))
			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
