package cli

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"mime/multipart"

	"github.com/spf13/cobra"
)




var sendCmd = &cobra.Command{
	Use:   "send [arquivo]",
	Short: "Send a file to a peer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		err := sendFile(filePath)
		if err != nil {
			fmt.Println("Error sending file:", err)
			os.Exit(1)
		}
	},
}


func init() {
	rootCmd.AddCommand(sendCmd)
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
		return fmt.Errorf("error sending file: %s", resp.Status)
	}

	return nil
}