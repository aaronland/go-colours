// Command line tool to extrude (derive) dominant colours from one or more images as well as closest matches colours
// using zero or more "snap-to-grid" colour palettes as JSON-encoded data written to STDOUT.
package extrude

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/aaronland/go-colours/extrude"
)

func Run(ctx context.Context) error {
	fs := DefaultFlagSet()
	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(fs)

	if err != nil {
		return err
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	if opts.Verbose {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		slog.Debug("Verbose logging enabled")
	}

	extrude_opts := &extrude.ExtrudeOptions{
		ExtruderURIs: opts.ExtruderURIs,
		PaletteURIs:  opts.PaletteURIs,
		Images:       opts.Images,
		AllowRemote:  opts.AllowRemote,
	}

	extrude_rsp, err := extrude.Extrude(ctx, extrude_opts)

	if err != nil {
		return fmt.Errorf("Failed to extrude images, %w", err)
	}

	if extrude_rsp.IsTmpRoot {
		defer os.RemoveAll(extrude_rsp.Root)
	}

	enc := json.NewEncoder(opts.Writer)
	err = enc.Encode(extrude_rsp.Images)

	if err != nil {
		return fmt.Errorf("Failed to encode images, %w", err)
	}

	return nil
}
