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
	X_high     float64 // Highest number in the x axis
	X_low      float64 // Lowest number in the x axis
	Y_high     float64 // Higest number in the y axis
	Y_low      float64 // Lowest number in the y axis
	X_start    int
	Y_start    int
	X_size     int
	Y_size     int
	Pixel_size float64
}

type Index struct {
	X int
	Y int
}
