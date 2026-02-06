package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCmdConfig() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "Manage Hermes configuration",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Config command called")
		},
	}
}
