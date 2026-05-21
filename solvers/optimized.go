package solvers

import (
	"fmt"
	"log/slog"
	"runtime"
	"sync"

	"github.com/navod-abay/mandelbrotset-go/models"
	"github.com/navod-abay/mandelbrotset-go/writers"
)

const (
	maximum_iteration_depth = 1000
)

func getHorizontalMiddleVerticalEdgePixel(pixelArray [][]models.ColorPixel, x int, y int, skip int) int {
	num := 0
	if pixelArray[x-skip][y].NumIterations == maximum_iteration_depth {
		num++
	}
	if pixelArray[x][y].NumIterations == maximum_iteration_depth {
		num++
	}
	if pixelArray[x+skip][y].NumIterations == maximum_iteration_depth {
		num++
	}
	return num
}

func getHorizontalEdgeVerticalMiddlePixel(pixelArray [][]models.ColorPixel, x int, y int, skip int) int {
	num := 0
	if pixelArray[x-skip/2][y].NumIterations == maximum_iteration_depth {
		num++
	}
	if pixelArray[x+skip/2][y].NumIterations == maximum_iteration_depth {
		num++
	}
	return num
}

func getHorizontalMiddleVerticalMiddlePixelMemberNeighbors(pixelArray [][]models.ColorPixel, x int, y int, skip int, rightColumn int, leftCol int) int {
	num := rightColumn + leftCol
	if pixelArray[x][y-skip/2].NumIterations == maximum_iteration_depth {
		num++
	}
	if pixelArray[x][y+skip/2].NumIterations == maximum_iteration_depth {
		num++
	}
	return num
}

func getHorizontalMiddleVerticalMiddlePixel(pixelArray [][]models.ColorPixel, x int, y int, skip int) int {
	num := 0
	if pixelArray[x][y+skip/2].NumIterations == maximum_iteration_depth {
		num++
	}
	if pixelArray[x][y].NumIterations == maximum_iteration_depth {
		num++
	}
	if pixelArray[x][y-skip/2].NumIterations == maximum_iteration_depth {
		num++
	}
	return num
}

func GetSubImageDimensionsArrays(imageDimensions models.ImageDimensions) []models.ImageDimensions {
	x_length := imageDimensions.X_size
	y_length := imageDimensions.Y_size
	slog.Debug("Starting subImageDimensions calculation", "x_length", x_length, "y_length", y_length)
	x_subdivisions := 1
	y_subdivisions := 1
	processorGroupSize := runtime.NumCPU()
	fmt.Printf("Number of processes: %v\n", processorGroupSize)
	for processorGroupSize > 1 {
		var longerAxisSubdivisions *int
		var longerAxisLength *int
		if x_length > y_length {
			longerAxisSubdivisions = &x_subdivisions
			longerAxisLength = &x_length
		} else {
			longerAxisSubdivisions = &y_subdivisions
			longerAxisLength = &y_length
		}
		if processorGroupSize%2 == 0 {
			*longerAxisSubdivisions *= 2
			processorGroupSize /= 2
			*longerAxisLength /= 2
		} else {
			break
		}
	}
	slog.Debug("Finished calculating subImage subdivisions", "x_subdivision", x_subdivisions, "y_subdivisions", y_subdivisions, "x_length", x_length, "y_length", y_length)
	x_pos := 0
	X_val := imageDimensions.X_low
	Y_val := imageDimensions.Y_low
	var newImageDimensions []models.ImageDimensions
	num := 0
	for i := 0; i < x_subdivisions; i++ {
		y_pos := 0
		for j := 0; j < y_subdivisions; j++ {
			fmt.Printf("\n\nProcuess Num: %v\n", num)
			newImageDimension := imageDimensions
			newImageDimension.X_start = x_pos
			newImageDimension.X_low = X_val
			fmt.Printf("X_Start: %v\n", newImageDimension.X_start)
			newImageDimension.Y_start = y_pos
			newImageDimension.Y_low = Y_val
			fmt.Printf("Y_Start: %v\n", newImageDimension.Y_start)
			y_pos += y_length
			Y_val += float64(y_length) * imageDimensions.Pixel_size
			newImageDimension.X_size = x_pos + x_length
			newImageDimension.X_high = X_val + float64(x_length)*imageDimensions.Pixel_size
			fmt.Printf("X_Size: %v\n", newImageDimension.X_size)
			newImageDimension.Y_size = y_pos
			newImageDimension.Y_high = Y_val
			fmt.Printf("Y_Size: %v\n", newImageDimension.Y_size)
			newImageDimensions = append(newImageDimensions, newImageDimension)
			num++
		}
		Y_val = imageDimensions.Y_low
		X_val += float64(x_length) * imageDimensions.Pixel_size
		x_pos += x_length
	}
	fmt.Printf(`Number of sub images: %v`, len(newImageDimensions))
	return newImageDimensions
}

func RunOneSkipPass(pixelArray [][]models.ColorPixel, imageDimensions models.ImageDimensions, skip int, saveSnapShotsFlag bool, waitGroup *sync.WaitGroup) {
	fmt.Printf("\n\nskip: %v\n", skip)
	slog.Debug("Starting calculating edges", "imageDimensions", imageDimensions)
	// Handling edges of the array without changing the maximum iterations according to the neighbors inclusion
	for j := imageDimensions.Y_start + skip/2; j < imageDimensions.Y_size; j += skip {
		pixelArray[imageDimensions.X_start][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_start][j].Number, maximum_iteration_depth)
		pixelArray[imageDimensions.X_size-1][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_size-1][j].Number, maximum_iteration_depth)
	}
	for i := imageDimensions.X_start + skip/2; i < imageDimensions.X_size; i += skip {
		pixelArray[i][imageDimensions.Y_start].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.Y_start].Number, maximum_iteration_depth)
		pixelArray[i][imageDimensions.Y_size-1].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.X_size-1].Number, maximum_iteration_depth)
	}

	slog.Debug("Finished handling edges for iteration, staring horizontal middle vertical edge", "skip", skip)
	// Handling horizontal middle vertical left edge pixels (8s in a 3 by 3 grid)
	var leftCol int
	var rightCol int
	var middleCol int
	// Handling the center of the arrays
	for i := imageDimensions.X_start + skip; i < imageDimensions.X_size-skip; i += skip {
		leftCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i, imageDimensions.Y_start, skip)
		rightCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i, imageDimensions.Y_start+skip, skip)
		for j := imageDimensions.Y_start + skip/2; j < imageDimensions.Y_size-2*skip; j += skip {
			pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (7-leftCol-rightCol)*100)
			leftCol = rightCol
			rightCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i, j+(3*skip/2), skip)
		}
	}

	slog.Debug("Starting Horizontal Edge Vertical Middle")
	// Handling Horizontal Edge Vertical Middle pixels (4/6s in a grid)
	for i := imageDimensions.X_start + skip/2; i < imageDimensions.X_size-skip; i += skip {
		leftCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, i, imageDimensions.Y_start, skip)
		middleCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, i, imageDimensions.Y_start+skip, skip)
		rightCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, i, imageDimensions.Y_start+2*skip, skip)
		for j := imageDimensions.Y_start + skip; j < imageDimensions.Y_size-skip; j += skip {
			pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (7-leftCol-rightCol-middleCol)*100)
			leftCol = middleCol
			middleCol = rightCol
			rightCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, i, j+skip, skip)
		}
	}

	// Handling Horizontal middle and Verticle Middle pixels
	slog.Debug("Starting Horizontal Middle Vertical Middle")
	for i := imageDimensions.X_start + skip/2; i < imageDimensions.X_size-skip; i += skip {
		for j := skip + imageDimensions.Y_start + skip/2; j < imageDimensions.Y_size-skip; j += skip {
			pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (9-getHorizontalMiddleVerticalMiddlePixelMemberNeighbors(pixelArray, i, j, skip, rightCol, leftCol))*100)
			leftCol = rightCol
			rightCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i, j+skip, skip)
		}
		leftCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+skip/2, imageDimensions.Y_start+skip/2, skip)
		rightCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+skip/2, imageDimensions.Y_start+skip, skip)
	}

	slog.Debug("Finished iteration with skip", "skip", skip)
	if saveSnapShotsFlag {
		waitGroup.Add(1)
		go writers.SaveCsvSnapshot(pixelArray, imageDimensions, skip, waitGroup)
	}
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
	slog.Debug("InitSkip", "init_skip", init_skip)
	for i := imageDimensions.X_start; i < imageDimensions.X_size; i++ {
		Y_val = imageDimensions.Y_low
		pixelArray[i] = make([]models.ColorPixel, imageDimensions.Y_size)
		for j := imageDimensions.Y_start; j < imageDimensions.Y_size; j++ {
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
	println("Finished initializing the array and calculating the initial pass ")
	if saveSnapShotsFlag {
		waitGroup.Add(1)
		go writers.SaveCsvSnapshot(pixelArray, imageDimensions, init_skip, &waitGroup)
	}

	for skip := init_skip; skip > 1; skip /= 2 {
		RunOneSkipPass(pixelArray, imageDimensions, skip, saveSnapShotsFlag, &waitGroup)
	}
	waitGroup.Wait()
	fmt.Println("Pixel Array initialization is over")

	return pixelArray
}

func SubImageOptimizedCalculation(imageDimensions models.ImageDimensions, pixelArray [][]models.ColorPixel, init_skip int, outerWaitGroup *sync.WaitGroup, saveSnapShotsFlag bool) [][]models.ColorPixel {
	defer outerWaitGroup.Done()
	var waitGroup sync.WaitGroup
	X_val := imageDimensions.X_low
	Y_val := imageDimensions.Y_low
	slog.Debug("InitSkip", "init_skip", init_skip)
	for i := 0; i < imageDimensions.X_size-imageDimensions.X_start; i++ {
		Y_val = imageDimensions.Y_low
		for j := 0; j < imageDimensions.Y_size-imageDimensions.Y_start; j++ {
			cur_num := complex(X_val, Y_val)
			pixelArray[i+imageDimensions.X_start][j+imageDimensions.Y_start].Number = cur_num
			if i%init_skip == 0 && j%init_skip == 0 {
				pixelArray[i+imageDimensions.X_start][j+imageDimensions.Y_start].NumIterations = checkMandelbrotSetInclusion(cur_num, maximum_iteration_depth)
				// slog.Debug("Checked complex num", "num", pixelArray[i][j].Number)
			}
			Y_val += imageDimensions.Pixel_size
		}
		X_val += imageDimensions.Pixel_size
	}
	fmt.Printf("Finished initializing the array and calculating the initial pass. saveSnapshotFlag %v ", saveSnapShotsFlag)
	if saveSnapShotsFlag {
		waitGroup.Add(1)
		go writers.SaveCsvSnapshot(pixelArray, imageDimensions, init_skip, &waitGroup)
	}
	for skip := init_skip; skip > 1; skip /= 2 {

		// Handling edges of the array without changing the maximum iterations according to the neighbors inclusion
		for j := imageDimensions.Y_start + skip/2; j < imageDimensions.Y_size; j += skip {
			pixelArray[imageDimensions.X_start][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_start][j].Number, maximum_iteration_depth)
			pixelArray[imageDimensions.X_size-1][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_size-1][j].Number, maximum_iteration_depth)
		}
		for i := imageDimensions.X_start / 2; i < imageDimensions.X_size; i += skip {
			pixelArray[i][imageDimensions.Y_start].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.Y_start].Number, maximum_iteration_depth)
			pixelArray[i][imageDimensions.Y_size-1].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.X_size-1].Number, maximum_iteration_depth)
		}
		slog.Debug("Finished handling edges for iteration", "skip", skip)

		// Handling horizontal middle vertical left edge pixels
		leftCol := getHorizontalMiddleVerticalEdgePixel(pixelArray, imageDimensions.X_start-skip/2, imageDimensions.Y_start+skip, skip)
		rightCol := getHorizontalMiddleVerticalEdgePixel(pixelArray, imageDimensions.X_start+skip/2, imageDimensions.Y_start+skip, skip)
		// Handling the center of the arrays
		for i := imageDimensions.X_start + skip/2; i < imageDimensions.X_size-skip; i += skip {
			for j := skip + imageDimensions.Y_start; j < imageDimensions.Y_size-skip; j += skip {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (6-leftCol-rightCol)*100)
				leftCol = rightCol
				rightCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i, j, skip)
			}
			leftCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i-skip/2, imageDimensions.Y_start+skip, skip)
			rightCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i+skip/2, imageDimensions.Y_start+skip, skip)
		}

		// Handling Horizontal Edge Vertical Middle pixels
		leftCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start, imageDimensions.Y_start+skip/2, skip)
		middleCol := getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+skip, imageDimensions.Y_start+skip/2, skip)
		rightCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+2*skip, imageDimensions.Y_start+skip/2, skip)

		for i := imageDimensions.X_start + skip; i < imageDimensions.X_size-skip; i += skip {
			for j := skip + imageDimensions.Y_start + skip/2; j < imageDimensions.Y_size-skip; j += skip {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (7-leftCol-rightCol-middleCol)*100)
				leftCol = middleCol
				middleCol = rightCol
				rightCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i, j+skip, skip)
			}
			leftCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start, imageDimensions.Y_start+skip/2, skip)
			middleCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+skip, imageDimensions.Y_start+skip/2, skip)
			rightCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+2*skip, imageDimensions.Y_start+skip/2, skip)
		}

		// Handling Horizontal middle and Verticle Middle pixels
		leftCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, imageDimensions.X_start+skip/2, imageDimensions.Y_start, skip)
		rightCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+2*skip, imageDimensions.Y_start+skip/2, skip)

		for i := imageDimensions.X_start + skip; i < imageDimensions.X_size-skip; i += skip {
			for j := skip + imageDimensions.Y_start + skip/2; j < imageDimensions.Y_size-skip; j += skip {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (9-getHorizontalMiddleVerticalMiddlePixelMemberNeighbors(pixelArray, i, j, skip, rightCol, leftCol))*100)
				leftCol = rightCol
				rightCol = getHorizontalMiddleVerticalEdgePixel(pixelArray, i, j+skip, skip)
			}
			leftCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+skip/2, imageDimensions.Y_start+skip/2, skip)
			rightCol = getHorizontalEdgeVerticalMiddlePixel(pixelArray, imageDimensions.X_start+skip/2, imageDimensions.Y_start+skip, skip)
		}

		slog.Debug("Finished iteration with skip", "skip", skip)
		if saveSnapShotsFlag {
			waitGroup.Add(1)
			go writers.SaveCsvSnapshot(pixelArray, imageDimensions, skip, &waitGroup)
		}
	}
	waitGroup.Wait()
	return pixelArray
}
