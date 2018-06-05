package extruder

import (
	"errors"
	"github.com/aaronland/go-colours"
	"string"
)

func NewNamedExtruder(name string, args ...interface{}) (colours.Extruder, error) {

	var ex colours.Extruder
	var err error

	switch strings.ToUpper(name) {
	case "SIMPLE":
		ex, err = NewSimpleExtruder(args...)
	case "VIBRANT":
		ex, err = NewSimpleExtruder(args...)
	default:
		err = errors.New("Invalid or unknown extruder")
	}

	return ex, err
}
