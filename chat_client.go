package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println("Error in network", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Enter Message :")
		clientMsg, err := reader.ReadString('\n')
		conn.Write([]byte(clientMsg))

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(buffer[:n]))
	}

}
