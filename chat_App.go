package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
)

const (
	DefaultPort     = "8080"
	DefaultProtocol = "tcp"
)

var (
	connections    []net.Conn
	users          = make(map[string]string)
	usersconnected = make(map[string]net.Conn)
)

func main() {

	//TODO:: Need to add the password hashing functionality

	// sample data for map
	users["anushka"] = "virat"
	port := flag.String("port", DefaultPort, "the port to use for the chat server")
	protocol := flag.String("protocol", DefaultProtocol, "the protocol to use for the chat server (tcp or udp)")
	flag.Parse()

	ln, err := net.Listen(*protocol, ":"+*port)
	if err != nil {
		fmt.Println("Error creating listener:", err)
		os.Exit(1)
	}
	defer ln.Close()

	fmt.Println("Chat server listening on port", *port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		connections = append(connections, conn)

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	fmt.Println("Client got connected")
	username, err := registerOrLogin(conn)
	if err != nil {
		fmt.Println("Error registering or logging in:", err)
		return
	}

	usersconnected[username] = conn
	fmt.Println(usersconnected)
	fmt.Println("username is ", username)

	notifyClients(fmt.Sprintf("%s has joined the chat.\n", username), username)

	chat(conn, username)
}

func registerOrLogin(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	fmt.Println("Went inside the login section")

	// Prompt the client to log in or register
	_, err := writer.WriteString("Enter 'login' to log in or 'register' to create a new account: ")
	if err != nil {
		return "", err
	}
	writer.Flush()

	// Read the client's response
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	response = strings.TrimSpace(response)

	// If the client chooses to register, prompt them to enter a new username and password
	if response == "register" {
		return register(conn, reader, writer)
	} else if response == "login" {
		return login(conn, reader, writer)
	} else {
		return "", fmt.Errorf("invalid response: %s", response)
	}
}

func register(conn net.Conn, reader *bufio.Reader, writer *bufio.Writer) (string, error) {
	_, err := writer.WriteString("Enter a new username: ")
	if err != nil {
		return "", err
	}
	writer.Flush()

	username, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	username = strings.TrimSpace(username)

	_, err = writer.WriteString("Enter a password: ")
	if err != nil {
		return "", err
	}
	writer.Flush()

	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	password = strings.TrimSpace(password)

	// Add the new user to the map
	users[username] = password
	fmt.Println("User insise register", users)

	_, err = writer.WriteString("Successfully registered. Please log in.\n")
	if err != nil {
		return "", err
	}
	writer.Flush()

	return login(conn, reader, writer)
}

func login(conn net.Conn, reader *bufio.Reader, writer *bufio.Writer) (string, error) {
	if err := prompt(writer, "Enter your username: "); err != nil {
		return "", err
	}

	username, err := readString(reader)
	if err != nil {
		return "", err
	}

	username = strings.TrimSpace(username)

	if err := prompt(writer, "Enter your password: "); err != nil {
		return "", err
	}

	password, err := readString(reader)
	if err != nil {
		return "", err
	}

	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return "", errors.New("username and password cannot be empty")
	}

	for savedusername, savedpassword := range users {
		if savedusername == username && savedpassword == password {
			if err := prompt(writer, "Login successful!\n"); err != nil {
				return "", err
			}
			return username, nil
		}
	}

	return username, nil
}

func prompt(writer *bufio.Writer, msg string) error {
	_, err := writer.WriteString(msg)
	if err != nil {
		return err
	}
	return writer.Flush()
}

func readString(reader *bufio.Reader) (string, error) {
	str, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(str), nil
}

func chat(conn net.Conn, username string) {

	fmt.Println("we went inside the chat function")

	reader := bufio.NewReader(conn)

	for {
		// Read a message from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		message = strings.TrimSpace(message)

		// If the message is a command, handle it
		if strings.HasPrefix(message, "/") {
			err = handleCommand(conn, message, username)
			if err != nil {
				fmt.Println("Error handling command:", err)
				break
			}
			continue
		}

		// Otherwise, send the message to all other clients
		notifyClients(fmt.Sprintf("%s: %s\n", username, message), username)
	}

	// Remove the connection from the slice of connections
	for i, c := range connections {
		if c == conn {
			connections = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	// Remove the user from the map of users
	delete(users, username)

	// Notify other clients that the user has left the chat
	notifyClients(fmt.Sprintf("%s has left the chat.\n", username), username)

	conn.Close()
}

func handleCommand(conn net.Conn, command string, username string) error {
	command = strings.ToLower(command)

	switch {
	case command == "/who":
		listUsers(conn)
	case strings.HasPrefix(command, "/msg "):
		sendPrivateMessage(conn, command, username)
	default:
		_, err := conn.Write([]byte(fmt.Sprintf("Invalid command: %s\n", command)))
		return err
	}
	return nil
}

func listUsers(conn net.Conn) {
	var usernames []string
	for username := range users {
		usernames = append(usernames, username)
	}
	sort.Strings(usernames)

	message := "Connected users:\n"
	for _, username := range usernames {
		message += fmt.Sprintf("- %s\n", username)
	}

	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func sendPrivateMessage(conn net.Conn, command string, username string) {
	parts := strings.SplitN(command, " ", 3)
	if len(parts) != 3 {
		_, err := conn.Write([]byte("Invalid syntax for private message\n"))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
		return
	}

	recipient := parts[1]
	message := parts[2]

	recipientConn, ok := usersconnected[recipient]
	if !ok {
		_, err := conn.Write([]byte(fmt.Sprintf("User %s is not connected\n", recipient)))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
		return
	}

	_, err := recipientConn.Write([]byte(fmt.Sprintf("(private) %s: %s\n", username, message)))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
}

func notifyClients(message string, sender string) {

	for _, conn := range connections {
		if usersconnected[sender] == conn {
			continue
		}

		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
	}
}
