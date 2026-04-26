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
		X_val += float64(x_length) * imageDimensions.Pixel_size
		x_pos += x_length
	}
	fmt.Printf(`Number of sub images: %v`, len(newImageDimensions))
	return newImageDimensions
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
	leftCol := 0
	middleCol := 0
	rightCol := 0
	for skip := init_skip / 2; skip >= 1; skip /= 2 {

		// Handling edges of the array without changing the maximum iterations according to the neighbors inclusion
		for j := imageDimensions.Y_start; j < imageDimensions.Y_size; j += skip {
			pixelArray[imageDimensions.X_start][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_start][j].Number, maximum_iteration_depth)
			pixelArray[imageDimensions.X_size-1][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_size-1][j].Number, 1000)
		}
		for i := imageDimensions.X_start; i < imageDimensions.X_size; i += skip {
			pixelArray[i][imageDimensions.Y_start].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.Y_start].Number, maximum_iteration_depth)
			pixelArray[i][imageDimensions.Y_size-1].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.X_size-1].Number, 1000)
		}
		slog.Debug("Finished handling edges for iteration", "skip", skip)

		middleCol = getNumInclusionsInColWithSkips(pixelArray, imageDimensions.X_start, skip, skip)
		rightCol = getNumInclusionsInColWithSkips(pixelArray, skip, skip, skip)
		// Handling the center of the arrays
		for i := skip; i < imageDimensions.X_size-skip; i += skip {
			for j := skip; j < imageDimensions.Y_size-skip; j += skip {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (10-leftCol-middleCol-rightCol)*100)

				leftCol = middleCol
				middleCol = rightCol
				rightCol = getNumInclusionsInColWithSkips(pixelArray, i, j, skip)
			}
		}
		slog.Debug("Finished iteration with skip", "skip", skip)
		waitGroup.Add(1)
		go writers.SaveCsvSnapshot(pixelArray, imageDimensions, skip, &waitGroup)
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
			pixelArray[i][j].Number = cur_num
			if i%init_skip == 0 && j%init_skip == 0 {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(cur_num, maximum_iteration_depth)
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
	leftCol := 0
	middleCol := 0
	rightCol := 0
	for skip := init_skip / 2; skip >= 1; skip /= 2 {

		// Handling edges of the array without changing the maximum iterations according to the neighbors inclusion
		for j := imageDimensions.Y_start; j < imageDimensions.Y_size; j += skip {
			pixelArray[imageDimensions.X_start][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_start][j].Number, maximum_iteration_depth)
			pixelArray[imageDimensions.X_size-1][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[imageDimensions.X_size-1][j].Number, 1000)
		}
		for i := imageDimensions.X_start; i < imageDimensions.X_size; i += skip {
			pixelArray[i][imageDimensions.Y_start].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.Y_start].Number, maximum_iteration_depth)
			pixelArray[i][imageDimensions.Y_size-1].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][imageDimensions.X_size-1].Number, 1000)
		}
		slog.Debug("Finished handling edges for iteration", "skip", skip)

		middleCol = getNumInclusionsInColWithSkips(pixelArray, imageDimensions.X_start, skip, skip)
		rightCol = getNumInclusionsInColWithSkips(pixelArray, skip, skip, skip)
		// Handling the center of the arrays
		for i := skip; i < imageDimensions.X_size-skip; i += skip {
			for j := skip; j < imageDimensions.Y_size-skip; j += skip {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(pixelArray[i][j].Number, (10-leftCol-middleCol-rightCol)*100)

				leftCol = middleCol
				middleCol = rightCol
				rightCol = getNumInclusionsInColWithSkips(pixelArray, i, j, skip)
			}
		}
		slog.Debug("Finished iteration with skip", "skip", skip)
		waitGroup.Add(1)
		go writers.SaveCsvSnapshot(pixelArray, imageDimensions, skip, &waitGroup)
	}
	waitGroup.Wait()
	return pixelArray
}
