package colours

type Colour interface {
     Name() string
     Reference() string
     Color image.Color
     Closest []Colour
}

type Palette interface {
     Reference() string
     Colours() []Colour
}

type Extruder interface {
     Colours(image.Image, int) ([]Colours, error)
}

type Grid interface {
     Closest(Colour, Palette) (Colours, error)
}