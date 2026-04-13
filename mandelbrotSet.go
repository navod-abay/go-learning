package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/navod-abay/mandelbrotset-go/models"
	"github.com/navod-abay/mandelbrotset-go/writers"
)

const (
	maximum_iteration_depth = 1000
	base_resolution         = 128
)

func getIntWithDefaultValue(prompt string, _default int) int {
	isSuccess := false
	var value int
	for !isSuccess {

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("%s [%d]:", prompt, _default)
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error Reading input.")
			fmt.Println(err)
			continue
		}
		str = strings.TrimSpace(str)
		if str == "" {
			value = _default
			isSuccess = true
		} else {
			var err error
			value, err = strconv.Atoi(str)
			if err != nil {
				fmt.Println("Error scanning the input")
				continue
			}
			isSuccess = true
		}
	}

	return value
}

func getFloatWithDefaultValue(prompt string, _default float64) float64 {
	isSuccess := false
	var value float64
	for !isSuccess {

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("%s [%f]:", prompt, _default)
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error Reading input.%v", err)
			continue
		}
		str = strings.TrimSpace(str)
		if str == "" {
			value = _default
			isSuccess = true
		} else {
			var err error
			value, err = strconv.ParseFloat(str, 64)
			if err != nil {
				fmt.Println("Error scanning the input")
				continue
			}
			isSuccess = true
		}
	}

	return value
}

func getInputWithCondition(prompt string, condition func(int) bool) int {
	input := getOneIntInput(prompt)
	for !condition(input) {
		fmt.Println("The input doesn't satisify the condition")
		getOneIntInput(prompt)
	}
	return input
}

func getOneIntInput(prompt string) int {
	isCorrect := false
	var i int
	for !isCorrect {
		fmt.Println(prompt)
		numScanned, err := fmt.Scan(&i)
		if numScanned > 1 {
			fmt.Println("Only 1 value is expected")
		} else if err != nil {
			fmt.Println("Error Scanning input")
		} else {
			return i
		}
	}
	return i
}

func calculatePixelSize(imageDimensions models.ImageDimensions, subdivision_level int) (float64, int, int, models.ImageDimensions) {
	X_range := imageDimensions.X_high - imageDimensions.X_low
	Y_range := imageDimensions.Y_high - imageDimensions.Y_low
	var pixel_size float64
	if X_range < Y_range {
		pixel_size = float64(X_range) / float64(base_resolution*(int(1)<<subdivision_level))
	} else {
		pixel_size = float64(Y_range) / float64(base_resolution*(int(1)<<subdivision_level))
	}
	X_size := int(float64(X_range) / pixel_size)
	Y_size := int(float64(Y_range) / pixel_size)
	updateImageDimensions := models.ImageDimensions{
		X_low:      imageDimensions.X_low,
		X_high:     imageDimensions.X_low + pixel_size*float64(X_size),
		Y_low:      imageDimensions.Y_low,
		Y_high:     imageDimensions.Y_high,
		Pixel_size: pixel_size,
		X_size:     X_size,
		Y_size:     Y_size,
	}
	return pixel_size, X_size, Y_size, updateImageDimensions
}

func complexMultiplicationSomponents(num complex128) (float64, float64, float64) {
	r := real(num)
	i := imag(num)
	return r * r, r * i, i * i
}

func checkMandelbrotSetInclusionNoColor(z_0 complex128, max_iteration int) bool {
	z_i := z_0
	r := real(z_0)
	i := imag(z_0)
	num := 0
	for ; num < max_iteration; num++ {
		rr, ri, ii := complexMultiplicationSomponents(z_i)
		if rr+ii > 4 {
			return false
		} else {
			z_i = complex(rr-ii+r, 2*ri+i)
		}
	}
	return true
}

func checkMandelbrotSetInclusion(z_0 complex128, max_iteration int) uint16 {
	z_i := z_0
	r := real(z_0)
	i := imag(z_0)
	var num uint16 = 0
	for ; num < uint16(max_iteration); num++ {
		rr, ri, ii := complexMultiplicationSomponents(z_i)
		if rr+ii > 4 {
			return num
		} else {
			z_i = complex(rr-ii+r, 2*ri+i)
		}
	}
	return maximum_iteration_depth
}

func ConstructAndCalculateNoColorPixelArray(imageDimensions models.ImageDimensions) [][]models.NoColorPixel {
	pixelArray := make([][]models.NoColorPixel, imageDimensions.X_size)
	for i := range pixelArray {
		pixelArray[i] = make([]models.NoColorPixel, imageDimensions.Y_size)
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
			pixelArray[i][j] = models.NoColorPixel{Number: cur_num, Included: checkMandelbrotSetInclusionNoColor(cur_num, maximum_iteration_depth)}
			y_val += imageDimensions.Pixel_size
		}
		x_val += imageDimensions.Pixel_size
	}
	return pixelArray
}

func ConstructAndCalculatePixelArray(imageDimensions models.ImageDimensions) [][]models.ColorPixel {
	pixelArray := make([][]models.ColorPixel, imageDimensions.X_size)
	for i := range pixelArray {
		pixelArray[i] = make([]models.ColorPixel, imageDimensions.Y_size)
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
			pixelArray[i][j] = models.ColorPixel{Number: cur_num, NumIterations: checkMandelbrotSetInclusion(cur_num, maximum_iteration_depth)}
			y_val += imageDimensions.Pixel_size
		}
		x_val += imageDimensions.Pixel_size
	}
	return pixelArray
}

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

func optimizedCalculation(imageDimensions models.ImageDimensions, subdivision_level int, saveSnapShotsFlag bool) [][]models.ColorPixel {
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
		pixelArray[i] = make([]models.ColorPixel, imageDimensions.Y_size)
		for j := range pixelArray[0] {
			cur_num := complex(X_val, Y_val)
			pixelArray[i][j].Number = cur_num
			if i%init_skip == 0 {
				pixelArray[i][j].NumIterations = checkMandelbrotSetInclusion(cur_num, maximum_iteration_depth)
			}
		}
	}
	println("Finished calculating the initial pass without optimization")
	if saveSnapShotsFlag {
		writers.SaveCsvSnapshot(pixelArray, imageDimensions, init_skip)
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
		writers.SaveCsvSnapshot(pixelArray, imageDimensions, skip)
	}
	fmt.Println("Pixel Array initialization is over")

	return pixelArray
}

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	var imageDimensions models.ImageDimensions
	imageDimensions.X_low = getFloatWithDefaultValue("Enter the x axis lower limit", -2)
	imageDimensions.X_high = getFloatWithDefaultValue("Enter the x axis upper limit", 2)
	imageDimensions.Y_low = getFloatWithDefaultValue("Enter y axis lower limit", -2)
	imageDimensions.Y_high = getFloatWithDefaultValue("Enter y axis upper limit", 2)
	subdivision_level := getIntWithDefaultValue("Enter the subdivision level", 4)
	var pixel_size float64
	var X_size, Y_size int
	pixel_size, X_size, Y_size, imageDimensions = calculatePixelSize(imageDimensions, subdivision_level)
	fmt.Println("Calculated array size")
	fmt.Printf("X axis size: %v\n", X_size)
	fmt.Printf("Y axis size: %v\n", Y_size)
	fmt.Printf("Pixel size: %v\n", pixel_size)

	optimizationFlag := flag.Bool("optimization", true, "Set to true to enable optimization")
	colorFlag := flag.Bool("not-colorized", true, "Produces a two colored image when set to true")
	csvWriteFlag := flag.Bool("write-csv", true, "Set to true to write the end result to a csv")
	bmpWriteFlag := flag.Bool("write-bmp", true, "Set to true to create a bmp image file")
	saveSnapShotsFlag := flag.Bool("save-snapshots", true, "Set to save intermediate results in the optimization process")

	flag.Parse()
	// Creating the image array
	if !*optimizationFlag {
		if *colorFlag {
			pixelArray := ConstructAndCalculatePixelArray(imageDimensions)
			if *csvWriteFlag {
				writers.WriteToCSV(pixelArray)
			}
			if *bmpWriteFlag {
				includedColor := make([]byte, 2)
				binary.LittleEndian.PutUint16(includedColor, 0)
				excludedColor := make([]byte, 2)
				binary.LittleEndian.PutUint16(excludedColor, 255)
				writers.WriteToBmpFile(pixelArray, imageDimensions, maximum_iteration_depth)

			}
		} else {
			pixelArray := ConstructAndCalculateNoColorPixelArray(imageDimensions)
			if *csvWriteFlag {
				writers.WriteToCSVNoColor(pixelArray)
			}
			if *bmpWriteFlag {
				includedColor := make([]byte, 2)
				binary.LittleEndian.PutUint16(includedColor, 0)
				excludedColor := make([]byte, 2)
				binary.LittleEndian.PutUint16(excludedColor, 255)
				writers.WriteToBmpFileNoColor(pixelArray, imageDimensions, includedColor, excludedColor)

			}
		}
	} else {
		pixelArray := optimizedCalculation(imageDimensions, subdivision_level, *saveSnapShotsFlag)
		if *csvWriteFlag {
			writers.WriteToCSV(pixelArray)
		}
		if *bmpWriteFlag {
			includedColor := make([]byte, 2)
			binary.LittleEndian.PutUint16(includedColor, 0)
			excludedColor := make([]byte, 2)
			binary.LittleEndian.PutUint16(excludedColor, 255)
			writers.WriteToBmpFile(pixelArray, imageDimensions, maximum_iteration_depth)

		}

	}
	fmt.Println("Calculation is over")

}
