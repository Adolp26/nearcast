package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error uploading file:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	nomeArquivo := header.Filename
	fmt.Println("Receiving file:", nomeArquivo, "size:", header.Size)

	dst, err := os.Create(nomeArquivo)
	if err != nil {
		fmt.Println("Error creating file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error saving file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Transfer finished:", nomeArquivo)
	fmt.Fprintf(w, "File %s uploaded successfully", nomeArquivo)
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Http server listening on port 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
