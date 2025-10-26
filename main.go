package main

import (
	"os"

	"github.com/MayR-Labs/mayrlabs-go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
