package cmd

import (
	"dnstt-tunnel/internal/server"

	"github.com/spf13/cobra"
)

var serverFlags = &server.Options{}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start DNS tunnel server",
	Run: func(cmd *cobra.Command, args []string) {
		s := server.New(serverFlags)
		s.Start()
	},
}

func init() {
	serverCmd.Flags().StringVarP(
		&serverFlags.HOST,
		"HOST",
		"H",
		"0.0.0.0",
		"server listen address",
	)

	serverCmd.Flags().IntVarP(
		&serverFlags.PORT,
		"PORT",
		"P",
		53,
		"server UDP port",
	)

	serverCmd.Flags().StringSliceVarP(
		&serverFlags.DOMAINS,
		"DOMAIN",
		"D",
		nil,
		"NS domain pointed to this server",
	)

	err := serverCmd.MarkFlagRequired("DOMAIN")
	if err != nil {
		panic(err)
	}

	rootCmd.AddCommand(serverCmd)
}
