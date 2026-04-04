package writers

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"

	"github.com/navod-abay/mandelbrotset-go/models"
)

type BmpHeaderDetails struct {
	fileSize       uint32
	reserved       uint32
	infoHeaderSize uint32
	dataOffset     uint32
	width          int32
	height         int32
	planes         uint16
	bitCount       uint16
	compression    int32
	imageSize      int32
	endInfoHeader  []int32
}

func WriteBmpHeader(file *os.File, headerDetails BmpHeaderDetails) {
	slog.Debug("Writing tp BMP file", "headerDetails", headerDetails)
	bufferedWriter := bufio.NewWriter(file)
	bufferedWriter.WriteString("BM")
	binary.Write(bufferedWriter, binary.LittleEndian, headerDetails.fileSize)
	binary.Write(bufferedWriter, binary.LittleEndian, headerDetails.reserved)
	binary.Write(bufferedWriter, binary.LittleEndian, headerDetails.dataOffset)
	binary.Write(bufferedWriter, binary.LittleEndian, headerDetails.infoHeaderSize)
	binary.Write(bufferedWriter, binary.LittleEndian, []int32{headerDetails.width, headerDetails.height})

	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, headerDetails.planes)
	bufferedWriter.Write(buf)

	binary.LittleEndian.PutUint16(buf, headerDetails.bitCount)
	bufferedWriter.Write(buf)

	binary.Write(bufferedWriter, binary.LittleEndian, headerDetails.endInfoHeader)
	bufferedWriter.Flush()
}

func CalculateBMPHeaderDetails(pixelArray [][]models.NoColorPixel) BmpHeaderDetails {
	var details BmpHeaderDetails
	details.infoHeaderSize = 40
	details.width = int32(len(pixelArray))
	details.height = int32(len(pixelArray[0]))
	details.planes = 1
	details.compression = 0
	details.imageSize = 0
	details.bitCount = 16
	details.dataOffset = 54
	details.reserved = 0
	details.fileSize = 2*uint32(details.width)*uint32(details.height) + 54
	fmt.Println("FileSize: ", details.fileSize)
	details.endInfoHeader = []int32{0, 0, 0, 0, 0, 0}
	return details
}

func WriteToBmpFileNoColor(pixelArray [][]models.NoColorPixel, includedColor []byte, excludedColor []byte) {
	fmt.Println("Writing output to bmp file")
	bmp_f, err := os.OpenFile("output.bmp", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Failed to opena  writer for the bmp file")
	} else {
		WriteBmpHeader(bmp_f, CalculateBMPHeaderDetails(pixelArray))
		writer := bufio.NewWriter(bmp_f)
		for i := range pixelArray[0] {
			for j := range pixelArray {
				if pixelArray[j][i].Included {
					writer.Write(includedColor)
				} else {
					writer.Write(excludedColor)
				}
			}
		}
		slog.Debug("Finished writing to the buffer")
		writer.Flush()
		slog.Debug("Flushed the buffer")
	}

	defer bmp_f.Close()
}
