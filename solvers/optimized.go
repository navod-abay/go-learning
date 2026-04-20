package solvers

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/navod-abay/mandelbrotset-go/models"
	"github.com/navod-abay/mandelbrotset-go/writers"
)

const (
	maximum_iteration_depth = 1000
)

// Calculate the number of pixels in the column is in the calculated to be in the set. Doesn't check for boundaries
func getNumInclusionsInColWithSkips(pixelArray [][]models.ColorPixel, x int, y int, skip int) int {
	num := 0
	if pixelArray[x][y-skip].NumIterations == maximum_iteration_depth {
		num++
	}
	if pixelArray[x][y].NumIterations == maximum_iteration_depth {
		num++
	}
	if pixelArray[x][y+skip].NumIterations == maximum_iteration_depth {
		num++
	}
	return num
}

func OptimizedCalculation(imageDimensions models.ImageDimensions, subdivision_level int, saveSnapShotsFlag bool) [][]models.ColorPixel {
	var waitGroup sync.WaitGroup
	pixelArray := make([][]models.ColorPixel, imageDimensions.X_size)
	var init_skip int
	if subdivision_level == 0 {
		init_skip = 1
	} else {
		init_skip = int(1) << (subdivision_level / 2)
	}
	X_val := imageDimensions.X_low
	Y_val := imageDimensions.Y_low
	slog.Debug("InitSkip", init_skip)
	for i := range pixelArray {
		Y_val = imageDimensions.Y_low
		pixelArray[i] = make([]models.ColorPixel, imageDimensions.Y_size)
		for j := range pixelArray[0] {
			cur_num := complex(X_val, Y_val)
			pixelArray[i][j].Number = cur_num
			if i%init_skip == 0 && j%init_skip == 0 {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(cur_num, maximum_iteration_depth)
				// slog.Debug("Checked complex num", "num", pixelArray[i][j].Number)
			}
			Y_val += imageDimensions.Pixel_size
		}
		X_val += imageDimensions.Pixel_size
	}
	println("Finished calculating the initial pass without optimization")
	if saveSnapShotsFlag {
		waitGroup.Add(1)
		go writers.SaveCsvSnapshot(pixelArray, imageDimensions, init_skip, &waitGroup)
	}
	leftCol := 0
	middleCol := 0
	rightCol := 0
	for skip := init_skip / 2; skip >= 1; skip /= 2 {

		// Handling edges of the array without changing the maximum iterations according to the neighbors inclusion
		for j := 0; j < imageDimensions.Y_size; j += skip {
			pixelArray[0][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[0][j].Number, maximum_iteration_depth)
			pixelArray[imageDimensions.X_size-1][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_size-1][j].Number, 1000)
		}
		for i := 0; i < imageDimensions.X_size; i += skip {
			pixelArray[i][0].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][0].Number, maximum_iteration_depth)
			pixelArray[0][imageDimensions.Y_size-1].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.X_size-1].Number, 1000)
		}
		slog.Debug("Finished handling edges for iteration", "skip", skip)

		middleCol = getNumInclusionsInColWithSkips(pixelArray, 0, skip, skip)
		rightCol = getNumInclusionsInColWithSkips(pixelArray, skip, skip, skip)
		// Handling the center of the arrays
		for i := skip; i < imageDimensions.X_size; i += skip {
			for j := skip; j < imageDimensions.Y_size; j += skip {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (10-leftCol-middleCol-rightCol)*100)

				leftCol = middleCol
				middleCol = rightCol
				rightCol = 0
			}
		}
		slog.Debug("Finished iteration with skip", "skip", skip)
		waitGroup.Add(1)
		writers.SaveCsvSnapshot(pixelArray, imageDimensions, skip, &waitGroup)
	}
	waitGroup.Wait()
	fmt.Println("Pixel Array initialization is over")

	return pixelArray
}
