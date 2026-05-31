package main

import (
	"fmt"
	"net"
	"sync"
)

const port string = ":8080"

func performHandshake(conn net.Conn) {
	message := "Hello from client 1"
	conn.Write([]byte(message))
}

func handleConnection(conn net.Conn, id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Connected with Id: %v", id)
	performHandshake(conn)
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
