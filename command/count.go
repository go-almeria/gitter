package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-almeria/gitx/api"
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

	args = flags.Args()
	if len(args) > 1 {
		flags.Usage()
		c.Ui.Error(fmt.Sprintf("\ncount expects at most one argument"))
		return 1
	}

	g := *api.NewGit("shortlog HEAD -n -s")
	g.Exec()
	lines, errc := g.Stream(os.Stdout)
	go func() {
		for {
			select {
			case line := <-lines:
				fmt.Println(line)
			case <-errc:
				return
			}
		}
	}()

	g.Wait()
	return 0
}

func (c *CountCommand) Synopsis() string {
	return "Outputs commit counts"
}

func (c *CountCommand) Help() string {
	helpText := `
Usage: gitx count [--all]
`
	return strings.TrimSpace(helpText)
}
