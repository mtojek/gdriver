package main

import (
	"os"

	"github.com/mtojek/gdriver/cmd"
)

func main() {
	rootCmd := cmd.Root()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
