package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"

	sharedproto "github.com/navod-abay/mandelbrotset-go/nodes/shared_proto"
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

func handleConnection(conn net.Conn, lock *sync.Mutex) error {
	lock.Lock() // only one master can push work into the client
	defer lock.Unlock()
	for true {
		bufReader := bufio.NewReader(conn)
		n, err := bufReader.Peek(1)
		if err != nil {
			continue // Blocking if the stream is empty
		}
		if len(n) > 0 {
			parse_err := parseMessage(bufReader)
			if parse_err != nil { // Close TCP connection if fail to parse
				return parse_err
			}
		}
	}
	return nil
}

func parseMessage(buffer *bufio.Reader) error {
	msgType, err := sharedproto.ReadHeader(buffer)
	if err != nil {
		switch {
		case errors.Is(err, sharedproto.ErrIncorrectMessage):

		}
	}
	return nil
}

func KeepListening(ln net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()
	var lock sync.Mutex
	fmt.Printf("Started Listening on port: %v\n", port)
	keepConnection := true
	for keepConnection {
		conn, err := ln.Accept() // Accept a connection
		if err != nil {
			fmt.Printf("Error accepting TCP message, err: %v ", err)
			continue
		}
		go handleConnection(conn, &lock)
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
