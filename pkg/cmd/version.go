package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cylonchau/hermes/pkg/version"
)

func NewCmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of hermes",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Hermes version: %s\n", version.Version)
		},
	}
}
