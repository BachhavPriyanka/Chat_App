package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strings"
)

const (
	DefaultPort     = "8080"
	DefaultProtocol = "tcp"
	DefaultServerIP = "127.0.0.1"
)

var (
	connections []net.Conn
)

func main() {
	port := flag.String("port", DefaultPort, "the port is used to chat server")
	protocol := flag.String("protocol", DefaultProtocol, "the protocol is used for chat server")
	serverip := flag.String("serverip", DefaultServerIP, "the address of server")
	flag.Parse()

	ln, err := net.Listen(*protocol, *serverip+":"+*port)
	//	ln, err := net.Listen("tcp ", ":8080")
	if err != nil {
		fmt.Println("Error in network", err)
	}
	defer ln.Close()

	fmt.Println("Server started. Listening on", *serverip+":"+*port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection", err)
			return
		}
		fmt.Println("New client connected:", conn.RemoteAddr())
		connections = append(connections, conn)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed:", conn.RemoteAddr())
			for i, c := range connections {
				if c == conn {
					connections = append(connections[:i], connections[i+1:]...)
					break
				}
			}
			return
		}

		message = strings.TrimSpace(message)

		for _, c := range connections {
			if c != conn {
				fmt.Fprintln(c, message)
			}
		}
	}

}
