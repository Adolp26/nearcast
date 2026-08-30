package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func handleConnection(conn net.Conn, id int) {
	defer conn.Close()

	tamanhoBuf := make([]byte, 4)

	_, err := io.ReadFull(conn, tamanhoBuf)
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return
	}

	tamanhoNome := binary.BigEndian.Uint32(tamanhoBuf)

	nomeBuf := make([]byte, tamanhoNome)

	_, err = io.ReadFull(conn, nomeBuf)
	if err != nil {
		fmt.Println("Error reading file name:", err)
		return
	}

	nomeArquivo := string(nomeBuf)

	fmt.Println("Receiving file:", nomeArquivo)

	tamanhoArquivoBuf := make([]byte, 4)
	_, err = io.ReadFull(conn, tamanhoArquivoBuf)
	if err != nil {
		fmt.Println("Error reading file size:", err)
		return
	}

	tamanhoArquivo := binary.BigEndian.Uint32(tamanhoArquivoBuf)

	file, err := os.Create(nomeArquivo)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = io.CopyN(file, conn, int64(tamanhoArquivo))
	if err != nil {
		fmt.Println("Error receiving file:", err)
		return
	}
	fmt.Println("Transfer finished for connection", id)
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("TCP server listening on port 9000")

	id := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Println("Connection accepted from", conn.RemoteAddr())

		id++
		go handleConnection(conn, id)
	}
}
