package review

import (
	"flag"
	"fmt"

	"github.com/aaronland/go-colours/extruder"
	"github.com/aaronland/go-colours/palette"
	"github.com/sfomuseum/go-flags/flagset"
)

type RunOptions struct {
	ExtruderURIs []string
	PaletteURIs  []string
	Root         string
	Images       []string
	Verbose      bool
}

func RunOptionsFromFlagSet(fs *flag.FlagSet) (*RunOptions, error) {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVars(fs, "REVIEW")

	if err != nil {
		return nil, fmt.Errorf("Failed to set flags from environment variables, %w", err)
	}

	if len(extruder_uris) == 0 {

		for _, scheme := range extruder.ExtruderSchemes() {
			extruder_uris.Set(scheme)
		}
	}

	if len(palette_uris) == 0 {

		for _, scheme := range palette.PaletteSchemes() {
			palette_uris.Set(scheme)
		}
	}

	opts := &RunOptions{
		ExtruderURIs: extruder_uris,
		PaletteURIs:  palette_uris,
		Root:         root,
		Images:       fs.Args(),
		Verbose:      verbose,
	}

	return opts, nil
}
