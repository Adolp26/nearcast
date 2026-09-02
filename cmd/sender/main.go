package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	err := sendFile("test.txt")
	if err != nil {
		fmt.Println("Error sending file:", err)
	}
}

func sendFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var corpo bytes.Buffer
	writer := multipart.NewWriter(&corpo)

	part, err := writer.CreateFormFile("file", "file.txt")
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:9000/upload", writer.FormDataContentType(), &corpo)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload file: %s", resp.StatusCode)
	}

	return nil
}
