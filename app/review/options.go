package review

import (
	"flag"
	"fmt"

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
		extruder_uris.Set("simple://")
		extruder_uris.Set("vibrant://")
	}

	if len(palette_uris) == 0 {
		palette_uris.Set("css3://")
		palette_uris.Set("css4://")
		palette_uris.Set("crayola://")
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
