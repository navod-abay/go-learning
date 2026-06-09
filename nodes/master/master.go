package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strconv"
	"sync"

	sharedproto "github.com/navod-abay/mandelbrotset-go/nodes/shared_proto"
)

type ClientIdentifier struct {
	ip           string
	numProcesses int
}

const port string = ":8080"

func sendHandshakeRequest(conn *bufio.Writer, id int) error {
	slog.Debug("Sending Handshake Request", "id:", id)
	err := sharedproto.SendMessage(conn, []byte{}, "handshake")
	slog.Debug("Handshake Request Sent")
	return err
}

func readHandshakeResponse(conn *bufio.Reader) (int, error) {
	requestType, bytes, err := sharedproto.ReadMessage(conn)
	if err != nil || requestType != sharedproto.HandshakeResponse {
		return 0, sharedproto.ErrIncorrectMessage
	}
	slog.Debug("Recieved Handshake Response", "received Extra content", bytes.String())
	// conn.Reset(conn)
	contentMap, err := sharedproto.ContentDeserialization(bytes)
	numProcesses, err := strconv.Atoi(contentMap["NumCPU"])
	return numProcesses, err
}

func handleConnection(conn net.Conn, id int, wg *sync.WaitGroup, c chan ClientIdentifier) {
	defer wg.Done()
	fmt.Printf("Connected with Id: %v\n", id)
	buffConn := bufio.NewWriter(conn)
	sendHandshakeRequest(buffConn, id)

	bufferedReader := bufio.NewReader(conn)
	slog.Debug("Gonna wait for the Handshake response")
	numProcesses, err := readHandshakeResponse(bufferedReader)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	clientIdentity := ClientIdentifier{ip: conn.LocalAddr().Network(), numProcesses: numProcesses}
	slog.Debug("Going to write the client identity to the channel", "clientIdentity", clientIdentity)
	c <- clientIdentity
	slog.Debug("Finished initiating connection", "clientIdentity", clientIdentity)
}

func main() {
	var opts = &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	var handler = slog.NewTextHandler(os.Stdin, opts)

	var logger = slog.New(handler)

	slog.SetDefault(logger)
	fmt.Printf("Running master node\n")
	IPs := []string{"127.0.0.1"}
	var wg sync.WaitGroup
	c := make(chan ClientIdentifier, 5)
	for id, ip := range IPs {
		fmt.Printf("Trying to connect to IP: %v\n", ip+port)
		conn, err := net.Dial("tcp", ip+port)
		if err != nil {
			fmt.Printf("Couldn't connect to IP: %v, err: %v\n", ip, err)
			continue
		}
		wg.Add(1)
		fmt.Printf("Successfully connected to ip: %v\n", ip)
		go handleConnection(conn, id, &wg, c)
	}
	wg.Wait()
	close(c)
	identities := map[string]ClientIdentifier{}
	for identity := range c {
		identities[identity.ip] = identity
	}
	fmt.Println("Ready slaves:", identities)
	fmt.Printf("Exiting master node\n")
}
