package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "dnstt-tunnel",
	Short: "DNS tunnel socks/server",
	Long:  "dnstt-tunnel is a DNS based VPN/tunnel",
}

func main() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}
