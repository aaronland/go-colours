package review

import (
	"flag"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/multi"
)

var extruder_uris multi.MultiString
var palette_uris multi.MultiString

var root string
var verbose bool

func DefaultFlagSet() *flag.FlagSet {

	fs := flagset.NewFlagSet("review")

	fs.Var(&extruder_uris, "extruder-uri", "Zero or more aaronland/go-colours/extruder.Extruder URIs. Default is to use all registered extruder schemes.")
	fs.Var(&palette_uris, "palette-uri", "Zero or more aaronland/go-colours/palette.Palette URIs. Default is to use all registered palette schemes.")
	fs.StringVar(&root, "root", "", "The path to a directory where images and HTML files associated with the review should be stored. If empty a new temporary directory will be created (and deleted when the application exits).")
	fs.BoolVar(&verbose, "verbose", false, "Enable verbose (debug) logging.")

	return fs
}
