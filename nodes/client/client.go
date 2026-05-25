package main

import (
	"fmt"
	"net"
	"sync"
)

const port string = ":8080"

func handleStringMessage(conn net.Conn, err error, keepListening *bool) (string, error) {
	buffer := make([]byte, 1024)
	byteCount, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Couldn't handle message in client, err: %v\n", err)
	}
	if byteCount > 0 {
		return string(buffer), nil
	}
	return "", err
}

func KeepListening(ln net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Started Listening on port: %v\n", port)
	keepConnection := true
	for keepConnection {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Error accepting TCP message, err: %v ", err)
			continue
		}
		handleStringMessage(conn, err, &keepConnection)
	}
}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	var wg = new(sync.WaitGroup)
	if err != nil {
		fmt.Printf("Can't connect to the client. Error: %v", err)
	}
	wg.Add(1)
	go KeepListening(ln, wg)
	wg.Wait()
}
