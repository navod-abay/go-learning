package models

type Pixel struct {
	Number   complex128
	Included byte
}

type ImageDimensions struct {
	X_high float64
	X_low  float64
	Y_high float64
	Y_low  float64
}
