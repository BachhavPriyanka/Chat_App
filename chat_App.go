package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in network", err)
	}
	defer conn.Close()

	for {
		ln, err := conn.Accept()
		if err != nil {
			fmt.Println("Error Accepting Connection", err)
			return
		}
		go handleConnection(ln)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	readerOS := bufio.NewReader(os.Stdin)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error in reading", err)
			return
		}
		fmt.Println("client -->", msg)

		fmt.Println("Enter message :")
		serverMsg, err := readerOS.ReadString('\n')
		if err != nil {
			fmt.Println("Error", err)
			return
		}
		conn.Write([]byte("Server --> " + serverMsg))
	}

}
