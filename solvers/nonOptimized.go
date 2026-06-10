package solvers

import (
	"fmt"

	"github.com/navod-abay/mandelbrotset-go/models"
)

func ConstructAndCalculateNoColorPixelArray(imageDimensions models.ImageDimensions) [][]bool {
	pixelArray := make([][]bool, imageDimensions.X_size)
	for i := range pixelArray {
		pixelArray[i] = make([]bool, imageDimensions.Y_size)
	}
	fmt.Println("Pixel Array initialization is over")

	// Populate the array

	// if no optimizatio is done, check the values while creation
	x_val := imageDimensions.X_low
	fmt.Println("Running in non optimized mode")
	for i := range pixelArray {
		y_val := imageDimensions.Y_low
		for j := range pixelArray[i] {
			cur_num := complex(x_val, y_val)
			pixelArray[i][j] = checkMandelbrotSetInclusionNoColor(cur_num, maximum_iteration_depth)
			y_val += imageDimensions.Pixel_size
		}
		x_val += imageDimensions.Pixel_size
	}
	return pixelArray
}

func ConstructAndCalculatePixelArray(imageDimensions models.ImageDimensions) [][]uint16 {
	pixelArray := make([][]uint16, imageDimensions.X_size)
	for i := range pixelArray {
		pixelArray[i] = make([]uint16, imageDimensions.Y_size)
	}
	fmt.Println("Pixel Array initialization is over")

	// Populate the array

	// if no optimizatio is done, check the values while creation
	x_val := imageDimensions.X_low
	fmt.Println("Running in non optimized mode")
	for i := range pixelArray {
		y_val := imageDimensions.Y_low
		for j := range pixelArray[i] {
			cur_num := complex(x_val, y_val)
			pixelArray[i][j] = checkMandelbrotSetInclusion(cur_num, maximum_iteration_depth)
			y_val += imageDimensions.Pixel_size
		}
		x_val += imageDimensions.Pixel_size
	}
	return pixelArray
}
