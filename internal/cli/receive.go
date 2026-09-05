package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const enderecoMulticast = "224.0.0.167:9999"

type AnuncioPeer struct {
	Alias string `json:"alias"` // as tags `json:"..."` dizem o nome do campo no JSON gerado
	Porta int    `json:"porta"`
}

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Inicia um servidor esperando arquivos",
	Run: func(cmd *cobra.Command, args []string) {
		go anunciarPresensa()
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

func anunciarPresensa() {

	addr, err := net.ResolveUDPAddr("udp", enderecoMulticast)
	if err != nil {
		fmt.Println("Error resolving multicast address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error dialing multicast address:", err)
		return
	}
	defer conn.Close()

	anuncio := AnuncioPeer{
		Alias: "Peer 1",
		Porta: 9000,
	}

	for {
		dados, err := json.Marshal(anuncio)
		if err != nil {
			fmt.Println("Error marshalling announcement:", err)
			return
		}

		_, err = conn.Write(dados)

		fmt.Println("Sent announcement:", string(dados))

		time.Sleep(3 * time.Second)
	}

}
