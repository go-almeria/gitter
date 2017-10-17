package meta

import "github.com/mitchellh/cli"

// Meta contains the meta-options and functionality that nearly every
// Gitx command inherits.
type Meta struct {
	Ui cli.Ui

	// These are set by the command line flags.
	flagClientKey string
}
