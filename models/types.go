package models

type NoColorPixel struct {
	Number   complex128
	Included bool
}

type ImageDimensions struct {
	X_high float64
	X_low  float64
	Y_high float64
	Y_low  float64
}
