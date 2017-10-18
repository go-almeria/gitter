package command

import (
	"strings"

	"github.com/go-almeria/gitx/meta"
)

// CountCommand Outputs commit count
type CountCommand struct {
	meta.Meta
}

func (c *CountCommand) Run(args []string) int {
	var all bool

	flags := c.Meta.FlagSet("count", meta.FlagSetDefault)
	flags.BoolVar(&all, "all", false, "")
	flags.Usage = func() { c.Ui.Error(c.Help()) }
	if err := flags.Parse(args); err != nil {
		return 1
	}

	return 0
}

func (c *CountCommand) Synopsis() string {
	return "Outputs commit count"
}

func (c *CountCommand) Help() string {
	helpText := `
Usage: gitx count [options]
`
	return strings.TrimSpace(helpText)
}
