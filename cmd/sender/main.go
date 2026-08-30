package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		fmt.Println("erro ao conectar:", err)
		return
	}
	defer conn.Close()

	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get file info:", err)
		return
	}

	tamanhoArquivo := info.Size()

	nomeArquivo := "test.txt"

	nomeArquivoBytes := []byte(nomeArquivo)

	tamanhoBuf := make([]byte, 4)

	tamanhoArquivoBuf := make([]byte, 4)

	binary.BigEndian.PutUint32(tamanhoBuf, uint32(len(nomeArquivoBytes)))
	binary.BigEndian.PutUint32(tamanhoArquivoBuf, uint32(tamanhoArquivo))

	conn.Write(tamanhoBuf)
	conn.Write(nomeArquivoBytes)
	conn.Write(tamanhoArquivoBuf)

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)

		if n > 0 {
			conn.Write(buffer[:n])
		}

		if err == io.EOF {
			fmt.Println("File sent successfully")
			break
		}
		if err != nil {
			fmt.Println("Failed to read file:", err)
			return
		}
	}
}
