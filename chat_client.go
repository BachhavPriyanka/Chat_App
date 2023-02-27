package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in network", err)
	}
	defer conn.Close()

	conn.Write([]byte("Hello Server"))

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(buffer[:n]))
}
