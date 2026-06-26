package sharedproto

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/navod-abay/mandelbrotset-go/core/models"
	"github.com/navod-abay/mandelbrotset-go/core/solvers"
)

type RpcServer int64

type EmptyInput struct{}

type StartWorkArgs struct {
	ImageDimensions   []models.ImageDimensions
	Subdivision_level int
}

func (server *RpcServer) HelloWorld(name string, out *string) error {
	fmt.Printf("Hello %v", name)
	*out = "Great to meet you"
	return nil
}

func (server *RpcServer) GetNumProcessor(nothing EmptyInput, numCPU *int) error {
	*numCPU = runtime.NumCPU()
	return nil
}

func (server *RpcServer) calculateSubImage(imageDimensions models.ImageDimensions, init_skip int32, c chan [][]uint16, outerWaitGroup *sync.WaitGroup) {
	pixelArray := make([][]uint16, imageDimensions.X_size)
	for i := range imageDimensions.X_size {
		pixelArray[i] = make([]uint16, imageDimensions.Y_size)
	}
	for i := range imageDimensions.X_size {
		pixelArray[i] = make([]uint16, imageDimensions.Y_size)
	}
	result := solvers.SubImageOptimizedCalculation(imageDimensions, pixelArray, init_skip, outerWaitGroup, false)
	c <- result
}

func (server *RpcServer) StartWork(startWorkArgs StartWorkArgs, result *int) error {
	subdivision_level := startWorkArgs.Subdivision_level
	fmt.Print("Running with parallelization")
	var waitGroup sync.WaitGroup // Wait group to wait for parallelized sub images
	var init_skip int32
	if subdivision_level == 0 {
		init_skip = 1
	} else {
		init_skip = int32(1) << (subdivision_level / 2)
	}
	subImageDimensionsArray := startWorkArgs.ImageDimensions
	c := make(chan [][]uint16)
	for _, subImageDimension := range subImageDimensionsArray {
		waitGroup.Add(1)
		go server.calculateSubImage(subImageDimension, init_skip, c, &waitGroup)
	}
	waitGroup.Wait()

	return nil
}
