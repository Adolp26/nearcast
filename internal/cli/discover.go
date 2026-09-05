package cli

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "Listen for available peers on the network",
	Run: func(cmd *cobra.Command, args []string) {
		descobrirPeers()
	},
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}

func descobrirPeers() {
	addr, err := net.ResolveUDPAddr("udp", enderecoMulticast)
	if err != nil {
		fmt.Println("Error resolving multicast address:", err)
		return
	}
	iface, err := net.InterfaceByName("Wi-Fi")
	if err != nil {
		fmt.Println("Error finding network interface:", err)
		return
	}

	conn, err := net.ListenMulticastUDP("udp", iface, addr)
	if err != nil {
		fmt.Println("Error listening multicast:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Listening peers...")

	buffer := make([]byte, 1024)

	for {
		n, remetente, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading UDP packet:", err)
			continue
		}

		var anuncio AnuncioPeer
		err = json.Unmarshal(buffer[:n], &anuncio)
		if err != nil {
			fmt.Println("Error decoding announcement:", err)
			continue
		}

		fmt.Printf("Peer discovered: %s at %s:%d\n", anuncio.Alias, remetente.IP, anuncio.Porta)
	}
}
