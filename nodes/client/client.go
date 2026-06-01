package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"sync"

	sharedproto "github.com/navod-abay/mandelbrotset-go/nodes/shared_proto"
)

const port string = ":8080"

func handleConnection(conn net.Conn, lock *sync.Mutex) error {
	lock.Lock() // only one master can push work into the client
	defer lock.Unlock()
	fmt.Println("Lock Acquired")
	for true {
		bufReader := bufio.NewReader(conn)
		n, err := bufReader.Peek(1)
		if err != nil {
			continue // Blocking if the stream is empty
		}
		if len(n) > 0 {
			msgType, parse_err := parseMessage(bufReader)
			if parse_err != nil { // Close TCP connection if fail to parse
				return parse_err
			}
			switch msgType {
			case sharedproto.HandShake:
				sendClientHandshakeResponse(conn)
			}
		}
	}
	return nil
}

func sendClientHandshakeResponse(conn net.Conn) {
	var buff bytes.Buffer
	sharedproto.WriteHeader(&buff, sharedproto.HandshakeResponse)
	cpus := runtime.NumCPU()
	buff.WriteString("\nNumCPU: " + strconv.Itoa(cpus) + "\n")
	conn.Write(buff.Bytes())

}

func parseMessage(buffer *bufio.Reader) (sharedproto.RequestType, error) {
	msgType, err := sharedproto.ReadHeader(buffer)
	if err != nil {
		switch {
		case errors.Is(err, sharedproto.ErrIncorrectMessage):

		}
	}
	fmt.Printf("message Type: %v", msgType)
	return msgType, nil
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
