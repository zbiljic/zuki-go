package main

import (
	"os"

	"github.com/zbiljic/zuki-go/examples/cobra/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
