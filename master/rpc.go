package main

import (
	"log/slog"
	"net/rpc"
	"sync"

	sharedproto "github.com/navod-abay/mandelbrotset-go/shared_proto"
)

func RpcHandshake(client *rpc.Client, wg *sync.WaitGroup, c chan ClientIdentifier) error {
	var numProcesses int16
	err := client.Call("RpcServer.GetNumProcessor", new(sharedproto.EmptyInput), &numProcesses)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	slog.Debug("NumProcess retrieval done", "numProcesses", numProcesses)
	return nil
}
