package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Value struct {
	typ   byte
	bytes []byte
	array []Value
}

// func decodeSimpleString(r *bufio.Reader) (Value, error) {
// 	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
// }
//
// func decodeBulkString(r *bufio.Reader) (Value, error) {
// 	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
// }
//
// func decodeArray(r *bufio.Reader) (Value, error) {
// 	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
// }

func DecodeRESP(byteStream *bufio.Reader) (Value, error) {
	dataTypeByte, err := byteStream.ReadByte()
	if err != nil {
		return Value{}, err
	}

	fmt.Println("dataTypeByte: ", dataTypeByte)
	fmt.Println("dataTypeByte: ", string(dataTypeByte))

	// switch string(dataTypeByte) {
	// case "+":
	// 	return decodeSimpleString(byteStream)
	// case "$":
	// 	return decodeBulkString(byteStream)
	// case "*":
	// 	return decodeArray(byteStream)
	// }

	return Value{}, nil
}

func handleConnect(conn net.Conn) {
	defer conn.Close()

	for {
		_, err := DecodeRESP(bufio.NewReader(conn))

		if err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			continue
		}

		conn.Write([]byte("+PONG\r\n"))
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		status, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			continue
		}
		fmt.Println("status: ", status)

		go handleConnect(conn)
	}
}
