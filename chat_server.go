package main

import (
	"fmt"
	"net"
)

func main() {

	fmt.Println("Starting TCP chat server")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in TCP server")
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		r, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error Reading", err)
			return
		}

		data := buffer[:r]
		fmt.Println("Data Recieved", string(data))

		conn.Write([]byte("Hello Client"))
	}
}
