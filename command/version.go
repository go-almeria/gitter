package command

import (
	"github.com/go-almeria/gitter/version"
	"github.com/mitchellh/cli"
)

// VersionCommand is a Command implementation prints the version.
type VersionCommand struct {
	VersionInfo *version.VersionInfo
	Ui          cli.Ui
}

func (c *VersionCommand) Help() string {
	return ""
}

func (c *VersionCommand) Run(_ []string) int {
	out := c.VersionInfo.FullVersionNumber(true)
	c.Ui.Output(out)
	return 0
}

func (c *VersionCommand) Synopsis() string {
	return "Prints the Gitter version"
}
