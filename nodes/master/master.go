package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"sync"

	sharedproto "github.com/navod-abay/mandelbrotset-go/nodes/shared_proto"
)

const port string = ":8080"

func handleConnection(conn net.Conn, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Connected with Id: %v", id)
	var buffer bytes.Buffer
	sharedproto.WriteHeader(&buffer, sharedproto.HandShake)
	conn.Write(buffer.Bytes())
	fmt.Println("Handshake Request Sent")
	buffer.Reset()
	io.Copy(conn, &buffer)
	bufferedReader := bufio.NewReader(conn)
	msgType, err := sharedproto.ReadHeader(bufferedReader)
	if err != nil {
		return
	}
	if msgType != sharedproto.HandshakeResponse {
		return
	}

}

func main() {
	fmt.Printf("Running master node\n")
	IPs := []string{"localhost", "127.0.0.1"}
	var wg sync.WaitGroup
	for id, ip := range IPs {
		fmt.Printf("Trying to connect to IP: %v\n", ip+port)
		conn, err := net.Dial("tcp", ip+port)
		if err != nil {
			fmt.Printf("Couldn't connect to IP: %v, err: %v\n", ip, err)
			continue
		}
		wg.Add(1)
		fmt.Printf("Successfully connected to ip: %v\n", ip)
		go handleConnection(conn, id, &wg)
	}
	wg.Wait()
	fmt.Printf("Exiting master node\n")
}
