package command

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
	"strconv"
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
	if err := g.Run(); err != nil {
		c.Ui.Error(fmt.Sprintf("%s", string(g.Err.Bytes())))
		return 1
	}

	re := regexp.MustCompile(`(\d+)\s+(\S+.*)`)
	scanner := bufio.NewScanner(&g.Out)

	count := 0
	for scanner.Scan() {
		results := re.FindAllStringSubmatch(scanner.Text(), -1)
		userCount, _ := strconv.Atoi(results[0][1])
		count += userCount
		user := results[0][2]
		if all {
			fmt.Printf("%s (%d)\n", user, userCount)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal("--->", err)
	}
	fmt.Printf("\ntotal %d\n", count)

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
