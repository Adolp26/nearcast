package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Inicia um servidor esperando arquivos",
	Run: func(cmd *cobra.Command, args []string) {
		startReceiver()
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
}

func startReceiver() {
	http.HandleFunc("/upload", uploadHandler)
	fmt.Println("Http server listening on port 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

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
		http.Error(w, "erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error saving file:", err)
		http.Error(w, "erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}

	fmt.Println("Transfer finished:", nomeArquivo)
	fmt.Fprintf(w, "File %s uploaded successfully", nomeArquivo)
}
