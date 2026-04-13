package models

type NoColorPixel struct {
	Number   complex128
	Included bool
}

type ColorPixel struct {
	Number        complex128
	NumIterations uint16
}

type ImageDimensions struct {
	X_high     float64
	X_low      float64
	Y_high     float64
	Y_low      float64
	X_size     int
	Y_size     int
	Pixel_size float64
}

type Index struct {
	X int
	Y int
}
