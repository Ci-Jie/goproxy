package cmd

import (
	"goproxy/server"

	"github.com/spf13/cobra"
)

func newStartCmd() (cmd *cobra.Command) {
	return &cobra.Command{
		Use:   "start",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start()
		},
	}
}
