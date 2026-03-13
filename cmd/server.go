package cmd

import (
	"dnstt-tunnel/server"
	"dnstt-tunnel/utils"

	"github.com/spf13/cobra"
)

var serverFlags = &server.Options{}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start DNS tunnel server",
	Run: func(cmd *cobra.Command, args []string) {
		if serverFlags.Password == "" {
			password, _ := utils.GeneratePassword(32)
			serverFlags.Password = password
		}

		s := server.New(serverFlags)
		s.Start()
	},
}

func init() {
	serverCmd.Flags().StringVarP(
		&serverFlags.Host,
		"host",
		"",
		"0.0.0.0",
		"UDP server host",
	)

	serverCmd.Flags().IntVarP(
		&serverFlags.Port,
		"port",
		"",
		53,
		"UDP server port",
	)

	serverCmd.Flags().StringVarP(
		&serverFlags.Domain,
		"domain",
		"d",
		"",
		"Domain to use for DNS query",
	)

	serverCmd.Flags().StringVarP(
		&serverFlags.Password,
		"password",
		"p",
		"",
		"Tunnel password",
	)

	serverCmd.Flags().IntVarP(
		&serverFlags.MaxRetransmits,
		"retries",
		"",
		5,
		"Max retransmits",
	)

	serverCmd.Flags().IntVarP(
		&serverFlags.FlowControlWindow,
		"window",
		"",
		4,
		"Flow control window size",
	)

	serverCmd.Flags().IntVarP(
		&serverFlags.KeepaliveInterval,
		"keepalive",
		"",
		5000,
		"Keepalive interval (ms)",
	)

	serverCmd.Flags().IntVarP(
		&serverFlags.AckTimeout,
		"ack-timeout",
		"",
		2000,
		"Wait time for ACK before retry (ms)",
	)

	serverCmd.Flags().IntVarP(
		&serverFlags.WriteTimeout,
		"write-timeout",
		"",
		2000,
		"Write timeout (ms)",
	)

	serverCmd.MarkFlagRequired("domain")
	rootCmd.AddCommand(serverCmd)
}
