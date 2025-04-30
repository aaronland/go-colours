package extruder

import (
	"image"

	"github.com/neocortical/noborders"
	"github.com/nfnt/resize"
)

func PrepareImage(im image.Image) (image.Image, error) {

	im = resize.Resize(100, 0, im, resize.Bilinear)

	opts := noborders.Opts()
	opts.SetEntropy(0.05)
	opts.SetVariance(100000)
	opts.SetMultiPass(true)

	return noborders.RemoveBorders(im, opts)
}
