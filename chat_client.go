package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	DefaultPort     = "8080"
	DefaultProtocol = "tcp"
	DefaultServerIp = "127.0.0.1"
)

func main() {
	port := flag.String("port", DefaultPort, "the port to use for chat server")
	protocol := flag.String("protocol", DefaultProtocol, "the protocol to use for chat server")
	serverip := flag.String("server", DefaultServerIp, "the address of server")

	flag.Parse()
	conn, err := net.Dial(*protocol, *serverip+":"+*port)
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
