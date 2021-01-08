package main

import (
	"os"

	"github.com/catenasys/sxtctl/pkg/commands"
)

func main() {
	if err := commands.NewCmdRoot(os.Stdout, os.Stderr).Execute(); err != nil {
		os.Exit(1)
	}
}
