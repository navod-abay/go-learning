package sharedproto

import (
	"fmt"
	"runtime"
)

type RpcServer int64

type EmptyInput struct{}

func (server *RpcServer) HelloWorld(name string, out *string) error {
	fmt.Printf("Hello %v", name)
	*out = "Great to meet you"
	return nil
}

func (serer *RpcServer) GetNumProcessor(nothing EmptyInput, numCPU *int) error {
	*numCPU = runtime.NumCPU()
	return nil
}
