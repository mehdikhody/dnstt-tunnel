package cmd

import (
	"dnstt-tunnel/internal/client"

	"github.com/spf13/cobra"
)

var clientFlags = &client.Options{}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start DNS tunnel socks",
	Run: func(cmd *cobra.Command, args []string) {
		s := client.New(clientFlags)
		s.Start()
	},
}

func init() {
	clientCmd.Flags().StringVarP(
		&clientFlags.HOST,
		"HOST",
		"H",
		"127.0.0.1",
		"socks5 server address",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.PORT,
		"PORT",
		"P",
		1080,
		"socks5 server port",
	)

	clientCmd.Flags().StringSliceVarP(
		&clientFlags.DOMAINS,
		"DOMAIN",
		"D",
		nil,
		"NS domain pointed to the server",
	)

	err := clientCmd.MarkFlagRequired("DOMAIN")
	if err != nil {
		panic(err)
	}

	clientCmd.Flags().Int64VarP(
		&clientFlags.MIN_UPLOAD_MTU,
		"MIN_UPLOAD_MTU",
		"",
		0,
		"Minimum upload chunk size per packet",
	)

	clientCmd.Flags().Int64VarP(
		&clientFlags.MAX_UPLOAD_MTU,
		"MAX_UPLOAD_MTU",
		"",
		1500,
		"Maximum upload chunk size per packet",
	)

	clientCmd.Flags().Int64VarP(
		&clientFlags.MIN_DOWNLOAD_MTU,
		"MIN_DOWNLOAD_MTU",
		"",
		0,
		"Minimum download chunk size per packet",
	)

	clientCmd.Flags().Int64VarP(
		&clientFlags.MAX_DOWNLOAD_MTU,
		"MAX_DOWNLOAD_MTU",
		"",
		1500,
		"Maximum download chunk size per packet",
	)

	clientCmd.Flags().IntVarP(
		&clientFlags.DNS_RESOLVER_TIMEOUT,
		"DNS_RESOLVER_TIMEOUT",
		"",
		3000,
		"DNS resolve timeout in milliseconds",
	)

	rootCmd.AddCommand(clientCmd)
}
