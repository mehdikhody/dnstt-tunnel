package cmd

import (
	"dnstt-tunnel/internal/server"
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate a private key",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := server.GenerateRandomKey()
		if err != nil {
			log.Printf("Failed to generate key: %v", err)
		}

		fmt.Printf("%x\n", key)
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)
}
