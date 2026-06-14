package sharedproto

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
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
	DelegateWork
)

var requestName = map[RequestType]string{
	HandShake:         "handshake",
	Invalid:           "invalid",
	HandshakeResponse: "handshakeResponse",
	DelegateWork:      "delegateWork",
}

var reverseMap = map[string]RequestType{
	"handshake":         HandShake,
	"invalid":           Invalid,
	"handshakeResponse": HandshakeResponse,
	"delegateWork":      DelegateWork,
}

func (rt RequestType) String() string {
	return requestName[rt]
}

func writeProtoIdenMsgType(buff *bytes.Buffer, requestType RequestType) error {
	headerString := ProtocolIdentifier + "\n" + requestName[requestType]
	n, err := buff.WriteString(headerString)
	if n < 1 || err != nil {
		return ErrNoBytesWritten
	}
	return nil
}

func checkHeader(buff *bufio.Reader) error {
	slog.Debug("Started checking the protocol identifier")
	line, err := buff.ReadString('\n')
	if err != nil {
		fmt.Printf("Error occured while reading the protocol identifier")
		switch {
		case errors.Is(err, io.EOF):
			return ErrStreamEmpty
		default:
			return err
		}
	}
	slog.Debug("Read string until newline")
	fmt.Printf("Identiied Protocol identifier: %v", line)
	if strings.TrimSpace(line) != ProtocolIdentifier {
		return ErrIncorrectMessage
	}
	fmt.Printf("Protocol identifier successfully read!!!!!")
	return nil
}

func readMessageType(buff *bufio.Reader) (RequestType, error) {
	line, err := buff.ReadString('\n')
	if err != nil {
		fmt.Println("Message type reading error")
	}
	fmt.Printf("Read messageType: %v", line)
	msgType, exists := reverseMap[strings.TrimSpace(line)]
	if exists {
		fmt.Printf("MessageType exists: %v", line)
		return msgType, nil
	}
	fmt.Printf("MessageType doesn't exists: %v", line)
	return Invalid, ErrIncorrectMessage
}

func ContentDeserialization(buff bytes.Buffer) (map[string]string, error) {
	slog.Debug("Started deserializing Content")
	content := make(map[string]string)
	line, err := buff.ReadString('\n')
	for err == nil {
		kvList := strings.Split(line, ":")
		if len(kvList) != 2 {
			slog.Error("Recieved a line with two delimiters", "kvList", kvList, "line", line)
			continue
		}
		content[kvList[0]] = strings.TrimSpace(kvList[1])
		line, err = buff.ReadString('\n')
	}
	slog.Debug("Finished Deserializing Content", "content", content)
	return content, nil
}

func ContentSerialization(buff *bytes.Buffer, content map[string]string) {
	for key, value := range content {
		buff.WriteString(key + ":" + value + "\n")
	}
}

func ReadNumProcesses(buff *bufio.Reader) (int, error) {
	fmt.Println("Started parsing num proceses")
	line, err := buff.ReadString('\n')
	if err != nil {
		return 0, ErrIncorrectMessage
	}
	lineContent := strings.Split(line, ": ")
	if len(lineContent) < 1 || lineContent[0] != "NumProcesses" || len(lineContent) < 2 {
		return 0, ErrIncorrectMessage
	}

	numProcesses, err := strconv.Atoi(lineContent[1])
	if err != nil {
		return 0, ErrIncorrectMessage
	}
	return numProcesses, nil
}

// Sends the message in buff along with the protocol identifier and the msgType Header
func SendMessage(conn *bufio.Writer, content []byte, msgType string) error {
	writeBuffer := bytes.Buffer{}
	err := writeProtoIdenMsgType(&writeBuffer, reverseMap[msgType])
	if err != nil {
		return err
	}
	writeBuffer.Write(content)
	writeBuffer.Write([]byte{'\n', '\n'})
	writeBuffer.WriteTo(conn)
	conn.Flush()
	return nil
}

// Reads the message from the connection and returns the msgType and the content bytes
func ReadMessage(conn *bufio.Reader) (RequestType, bytes.Buffer, error) {
	err := checkHeader(conn)
	if err != nil {
		return Invalid, bytes.Buffer{}, ErrIncorrectMessage
	}
	msgType, err := readMessageType(conn)
	if err != nil {
		return Invalid, bytes.Buffer{}, ErrIncorrectMessage
	}
	buffer := bytes.Buffer{}
	nextByte, err := conn.Peek(1)
	for nextByte[0] != '\n' && err == nil {
		byteArray, err := conn.ReadBytes('\n')
		if err != nil {
			break
		}
		buffer.Write(byteArray)
		nextByte, err = conn.Peek(1)
	}

	return msgType, buffer, nil
}
