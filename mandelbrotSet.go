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
	"github.com/navod-abay/mandelbrotset-go/solvers"
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
			pixelArray := solvers.ConstructAndCalculatePixelArray(imageDimensions)
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
			pixelArray := solvers.ConstructAndCalculateNoColorPixelArray(imageDimensions)
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
		pixelArray := solvers.OptimizedCalculation(imageDimensions, subdivision_level, *saveSnapShotsFlag)
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
