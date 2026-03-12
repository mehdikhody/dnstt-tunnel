package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dnstt-tunnel",
	Short: "DNS tunnel socks/server",
	Long:  "dnstt-tunnel is a DNS based VPN/tunnel",
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
