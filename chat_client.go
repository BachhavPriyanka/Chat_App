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
	DefaultServerIP = "127.0.0.1"
)

func main() {
	port := flag.String("port", DefaultPort, "the port to use for chat server")
	protocol := flag.String("protocol", DefaultProtocol, "the protocol to use for chat server")
	serverip := flag.String("serverip", DefaultServerIP, "the address of server")
	flag.Parse()

	conn, err := net.Dial(*protocol, *serverip+":"+*port)
	if err != nil {
		fmt.Println("Error in network", err)
	}
	defer conn.Close()

	msgChan := make(chan string)
	clientMsgChan := make(chan string)

	go func() {
		recivedserverbuffer := make([]byte, 1024)
		for {
			n, err := conn.Read(recivedserverbuffer)
			if err != nil {
				fmt.Println("Error receiving message:", err)
				continue
			}
			msgChan <- string(recivedserverbuffer[:n])
		}
	}()

	go func() {
		cmdDataReader := bufio.NewReader(os.Stdin)
		for {
			clientMsg, err := cmdDataReader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				continue
			}
			clientMsgChan <- clientMsg
		}
	}()

	for {
		fmt.Print("Enter message: ")
		fmt.Println()
		select {
		case msg := <-msgChan:
			fmt.Println("Received message:", msg)
		case clientMsg := <-clientMsgChan:
			_, err = conn.Write([]byte(clientMsg))
			if err != nil {
				fmt.Println("Error sending message:", err)
				continue
			}
		}
	}
}
