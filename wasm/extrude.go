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

	"github.com/aaronland/go-colours/grid"
	"github.com/aaronland/go-colours/palette"
	"github.com/aaronland/go-colours/extruder"			
	"github.com/aaronland/go-colours/extrude"
)

type ExtrudeOptions struct {
	Grid string `json:"grid"`
	Palettes []string `json:"palettes"`
	Extruders []string `json:"extruders"`
}

func (o *ExtrudeOptions) DeriveExtrusionOptions(ctx context.Context) (*extrude.DeriveExtrusionsOptions, error) {

	gr, err := grid.NewGrid(ctx, o.Grid)

	if err != nil {
		return nil, err
	}

	palettes := make([]palette.Palette, len(o.Palettes))
	extruders := make([]extruder.Extruder, len(o.Extruders))
	
	for idx, uri := range o.Palettes {
		
		p, err := palette.NewPalette(ctx, uri)

		if err != nil {
			return nil, err
		}

		palettes[idx] = p
	}

	for idx, uri := range o.Extruders {
		
		p, err := extruder.NewExtruder(ctx, uri)

		if err != nil {
			return nil, err
		}

		extruders[idx] = p
	}

	derive_opts := &extrude.DeriveExtrusionsOptions{
		Grid: gr,
		Palettes: palettes,
		Extruders: extruders,
	}

	return derive_opts, nil
}

func ExtrudeFunc () js.Func {

	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		str_opts := args[0].String()
		im_b64 := args[1].String()

		handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {

			resolve := args[0]
			reject := args[1]

			ctx := context.Background()
			
			var extrude_opts *ExtrudeOptions

			err := json.Unmarshal([]byte(str_opts), &extrude_opts)

			if err != nil {
				reject.Invoke(fmt.Errorf("Failed to decode derive options, %w", err))
				return nil
			}

			derive_opts, err := extrude_opts.DeriveExtrusionOptions(ctx)

			if err != nil {
				reject.Invoke(fmt.Errorf("Failed to derive extrusion options, %w", err))
				return nil
			}
			im_data, err := base64.StdEncoding.DecodeString(im_b64)
			
			if err != nil {
				reject.Invoke(fmt.Errorf("Failed to decode image body, %w", err))
				return nil
			}

			im_r := bytes.NewReader(im_data)
			
			im, _, err := image.Decode(im_r)

			if err != nil {
				reject.Invoke(fmt.Errorf("Failed to decode image, %w", err))
				return nil
			}
			
			rsp, err := extrude.DeriveExtrusions(ctx, derive_opts, im)

			if err != nil {
				reject.Invoke(fmt.Errorf("Failed to extrude image colours, %w", err))
				return nil
			}
			
			enc_rsp, err := json.Marshal(rsp)

			if err != nil {
				reject.Invoke(fmt.Errorf("Failed to marshal image colours, %w", err))
				return nil
			}
			
			resolve.Invoke(string(enc_rsp))
			return nil
		})

		promiseConstructor := js.Global().Get("Promise")
		return promiseConstructor.New(handler)
	})
}
