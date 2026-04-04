package writers

import (
	"bufio"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/navod-abay/mandelbrotset-go/models"
)

func WriteToCSVNoColor(pixelArray [][]models.NoColorPixel) {

	fmt.Println("Writing output to a csv file(No Colors)")
	f, err := os.OpenFile("outputNoColor.csv", os.O_WRONLY|os.O_CREATE, 0644)
	writer := bufio.NewWriter(f)

	if err == nil {
		for i := range pixelArray {
			for j := range pixelArray[i] {
				if pixelArray[i][j].Included {
					writer.WriteString("1, ")
				} else {
					writer.WriteString("0,")
				}
			}
			writer.WriteString("\n")
		}
		slog.Debug("Finished writing to the buffer")
		writer.Flush()
		slog.Debug("Flushed the buffer")
	} else {
		log.Fatal(err)
	}
	defer f.Close()
}

func WriteToCSV(pixelArray [][]models.ColorPixel) {

	fmt.Println("Writing output to a csv file")
	f, err := os.OpenFile("output.csv", os.O_WRONLY|os.O_CREATE, 0644)
	writer := bufio.NewWriter(f)

	if err == nil {
		for i := range pixelArray {
			for j := range pixelArray[i] {
				writer.WriteString(string(int32(pixelArray[i][j].NumIterations)))
			}
			writer.WriteString("\n")
		}
		slog.Debug("Finished writing to the buffer")
		writer.Flush()
		slog.Debug("Flushed the buffer")
	} else {
		log.Fatal(err)
	}
	defer f.Close()
}
