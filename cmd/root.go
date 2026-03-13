package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "dnstt-tunnel",
	Short: "DNS tunnel socks/server",
	Long:  "dnstt-tunnel is a DNS based VPN/tunnel",
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}
