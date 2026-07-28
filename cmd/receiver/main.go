package main

import (
	"fmt"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port 9000")

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connection accepted from", conn.RemoteAddr())

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading message:", err)
		return
	}

	fmt.Println("Received message:", string(buffer[:n]))}
