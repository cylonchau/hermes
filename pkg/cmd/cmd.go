package cmd

import (
	"github.com/cylonchau/hermes/pkg/cmd/server"
	"github.com/spf13/cobra"
)

type HermesOptions struct{}

func NewHermesCommand(o HermesOptions) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "hermes",
		Short: "hermes controls the Hermes DNS governance platform",
		Long: `
      Hermes controls the Hermes DNS governance platform.

      Find more information at:
            https://github.com/cylonchau/hermes/doc`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(
		NewCmdVersion(),
		server.NewCommand(),
	)
	return rootCmd
}
