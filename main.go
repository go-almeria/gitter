package main

import (
	"os"

	"github.com/go-almeria/gitx/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:]))
}
