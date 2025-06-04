//go:build wasmjs
package wasm

import (
	"context"
	"fmt"
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"syscall/js"
	"encoding/json"
	"encoding/base64"
	
	"github.com/aaronland/go-colours/extrude"
)

func ExtrudeFunc () js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		str_opts := args[0].String()
		im_b64 := args[1].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			ctx := context.Background()
			
			var derive_opts *extrude.DeriveExtrusionsOptions

			err := json.Unmarshal([]byte(str_opts), &derive_opts)

			if err != nil {
				reject.Invoke(fmt.Sprintf("Failed to decode derive options, %w", err))
				return nil
			}
			
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
			
			rsp, err := extrude.DeriveExtrusions(ctx, derive_opts, im)

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
