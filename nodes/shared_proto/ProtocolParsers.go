package sharedproto

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

const ProtocolIdentifier = "M@d3lbr0t$3t"

var (
	ErrStreamEmpty      = errors.New("The stream is Empty")
	ErrIncorrectMessage = errors.New("Incorrect Message")
)

type RequestType int

const (
	Invalid RequestType = iota
	HandShake
)

var requestName = map[RequestType]string{
	HandShake: "handshake",
	Invalid:   "invalid",
}

var reverseMap = map[string]RequestType{
	"handshake": HandShake,
	"invalid":   Invalid,
}

func (rt RequestType) String() string {
	return requestName[rt]
}

func ReadHeader(buff *bufio.Reader) (RequestType, error) {
	line, err := buff.ReadString('\n')
	if err != nil {
		switch {
		case errors.Is(err, io.EOF):
			return Invalid, ErrStreamEmpty
		default:
			return Invalid, err
		}
	}
	if strings.TrimSpace(line) != ProtocolIdentifier {
		return Invalid, ErrIncorrectMessage
	}
	line, err = buff.ReadString('\n')
	msgType, exists := reverseMap[strings.TrimSpace(line)]
	if exists {
		return msgType, nil
	}
	return Invalid, ErrIncorrectMessage
}
