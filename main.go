package main

import (
	"os"

	"github.com/go-almeria/gitter/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
