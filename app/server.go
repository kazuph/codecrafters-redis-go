package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		log.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	mem := NewMem()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnect(conn, mem)
	}
}

func handleConnect(conn net.Conn, mem *Mem) {
	defer conn.Close()

	for {
		msg := make([]byte, 1024)
		_, err := conn.Read(msg)
		if err != nil {
			log.Println("Error reading: ", err.Error())
			if errors.Is(err, io.EOF) {
				break // 重要
			}
			log.Fatal(err)
		}

		log.Println("msg: ", string(msg))

		ok := strings.HasSuffix(string(msg), "\r\n")
		if !ok {
			log.Println("There is no line break at the end.")
			msg = []byte(string(msg) + "\r\n")
		}

		value, err := DecodeRESP(bufio.NewReader(strings.NewReader(string(msg))))
		// value, err := DecodeRESP(bufio.NewReader(conn))

		if err != nil {
			log.Println("Error reading from client: ", err.Error())
			continue
		}

		if len(value.Array()) == 0 {
			conn.Write([]byte("-ERR unknown command\r\n"))
			panic("unknown command")
		}
		command := value.Array()[0].String()
		args := value.Array()[1:]

		switch command {
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case "echo":
			conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].String()), args[0].String())))
		case "set":
			log.Printf("SET key: %s, value: %s\n", args[0].String(), args[1].String())
			// len > 2の場合はオプションが存在する
			if len(args) > 2 {
				option := args[2].String()
				switch option {
				case "px":
					// pxの場合はミリ秒
					log.Printf("px: %s\n", args[3].String())
					// to int
					expireMSec, err := strconv.Atoi(args[3].String())
					if err != nil {
						conn.Write([]byte("-ERR value is not an integer or out of range\r\n"))
					}
					mem.SetWithExpiry(
						args[0].String(),
						args[1].String(),
						time.Duration(expireMSec)*time.Millisecond,
					)
				}

			} else {
				mem.Set(args[0].String(), args[1].String())
			}
			conn.Write([]byte("+OK\r\n"))
		case "get":
			key := args[0].String()
			value, found := mem.Get(key)
			log.Printf("GET key: %s, value: %s\n", key, value)
			if found {
				conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)))
			} else {
				conn.Write([]byte("$-1\r\n"))
			}
		default:
			conn.Write([]byte("-ERR unknown command '" + command + "'\r\n"))
		}

	}
}
