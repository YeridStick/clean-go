package main

import (
	"os"

	"github.com/YeridStick/cleango/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
