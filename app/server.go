package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	mem := NewMem()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnect(conn, mem)
	}
}

func handleConnect(conn net.Conn, mem *Mem) {
	defer conn.Close()

	for {
		value, err := DecodeRESP(bufio.NewReader(conn))

		if err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			continue
		}

		command := value.Array()[0].String()
		var args []Value
		if len(value.Array()) > 0 {
			args = value.Array()[1:]
		}

		switch command {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
		case "set":
			mem.Set(args[0].String(), args[1].String())
			conn.Write([]byte("+OK\r\n"))
		case "get":
			key := args[0].String()
			value := mem.Get(key)
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}

	}
}
