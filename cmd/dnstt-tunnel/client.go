package cmd

import (
	"dnstt-tunnel/client"
	"fmt"

	"github.com/spf13/cobra"
)

var clientFlags = &client.Options{}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start DNS tunnel socks",
	Run: func(cmd *cobra.Command, args []string) {
		hasUser := clientFlags.SocksUsername != ""
		hasPass := clientFlags.SocksPassword != ""
		if hasUser && !hasPass {
			fmt.Println("You need to specify --socks-password")
			return
		}
		if !hasUser && hasPass {
			fmt.Println("You need to specify --socks-username")
			return
		}

		s := client.New(clientFlags)
		s.Start()
	},
}

func init() {
	clientCmd.Flags().StringVarP(
		&clientFlags.SocksHost,
		"socks-host",
		"",
		"0.0.0.0",
		"Socks5 server host",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.SocksPort,
		"socks-port",
		"",
		1080,
		"Socks5 server port",
	)

	clientCmd.Flags().StringVarP(
		&clientFlags.SocksUsername,
		"socks-username",
		"",
		"",
		"Socks5 username",
	)

	clientCmd.Flags().StringVarP(
		&clientFlags.SocksPassword,
		"socks-password",
		"",
		"",
		"Socks5 password",
	)

	clientCmd.Flags().StringVarP(
		&clientFlags.Domain,
		"domain",
		"d",
		"",
		"Domain to use for DNS query",
	)

	clientCmd.Flags().StringVarP(
		&clientFlags.Password,
		"password",
		"p",
		"",
		"Tunnel password",
	)

	clientCmd.Flags().StringSliceVarP(
		&clientFlags.Resolvers,
		"resolver",
		"r",
		[]string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "1.0.0.1"},
		"Tunnel password",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.ChunkSize,
		"chunk-size",
		"",
		32,
		"DNS payload fragment size",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.MaxRetransmits,
		"retries",
		"",
		5,
		"Max retransmits",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.FlowControlWindow,
		"window",
		"",
		4,
		"Flow control window size",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.KeepaliveInterval,
		"keepalive",
		"",
		5000,
		"Keepalive interval (ms)",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.ReadPollInterval,
		"poll",
		"",
		400,
		"Poll interval for server data (ms)",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.AckTimeout,
		"ack-timeout",
		"",
		2000,
		"Wait time for ACK before retry (ms)",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.WriteTimeout,
		"write-timeout",
		"",
		2000,
		"Write timeout (ms)",
	)

	clientCmd.MarkFlagRequired("domain")
	clientCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(clientCmd)
}
