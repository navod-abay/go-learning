package sharedproto

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const ProtocolIdentifier = "M@nd3lbr0t$3t"

var (
	ErrStreamEmpty      = errors.New("The stream is Empty")
	ErrIncorrectMessage = errors.New("Incorrect Message")
	ErrNoBytesWritten   = errors.New("No Bytes Written")
)

type RequestType int

const (
	Invalid RequestType = iota
	HandShake
	HandshakeResponse
)

var requestName = map[RequestType]string{
	HandShake:         "handshake",
	Invalid:           "invalid",
	HandshakeResponse: "handshakeReponse",
}

var reverseMap = map[string]RequestType{
	"handshake":        HandShake,
	"invalid":          Invalid,
	"handshakeReponse": HandshakeResponse,
}

func (rt RequestType) String() string {
	return requestName[rt]
}

func WriteHeader(buff *bytes.Buffer, requestType RequestType) error {
	headerString := ProtocolIdentifier + "\n" + requestName[requestType]
	n, err := buff.WriteString(headerString)
	if n < 1 || err != nil {
		return ErrNoBytesWritten
	}
	return nil
}

func ReadHeader(buff *bufio.Reader) (RequestType, error) {
	fmt.Println("Started parsing the headers")
	line, err := buff.ReadString('\n')
	if err != nil {
		fmt.Printf("Error occured while reading the protocol identifier")
		switch {
		case errors.Is(err, io.EOF):
			return Invalid, ErrStreamEmpty
		default:
			return Invalid, err
		}
	}
	fmt.Printf("Identiied Protocol identifier: %v", line)
	if strings.TrimSpace(line) != ProtocolIdentifier {
		return Invalid, ErrIncorrectMessage
	}
	fmt.Printf("Protocol identifier successfully read")
	line, err = buff.ReadString('\n')
	fmt.Printf("Read messageType: %v", line)
	msgType, exists := reverseMap[strings.TrimSpace(line)]
	if exists {
		return msgType, nil
	}
	return Invalid, ErrIncorrectMessage
}
