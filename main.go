package main

import (
	"flag"
	"os"

	"github.com/cylonchau/hermes/pkg/cmd"
	"github.com/spf13/pflag"
)

func main() {
	command := cmd.NewHermesCommand(cmd.HermesOptions{})
	flagset := flag.CommandLine
	pflag.CommandLine.AddGoFlagSet(flagset)
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
