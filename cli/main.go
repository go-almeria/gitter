package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/mitchellh/cli"
)

func Run(args []string) int {
	return RunCustom(args, Commands(nil))
}

func RunCustom(args []string, commands map[string]cli.CommandFactory) int {
	_, err := exec.LookPath("git")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Git executable not reachable: %s\n", err.Error())
		return 1
	}

	cmd := exec.Command("git", "--version")
	_, err = cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Git executable not usable: %s\n", err.Error())
		return 1
	}

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	for _, arg := range args {
		if arg == "-v" || arg == "-version" || arg == "--version" {
			newArgs := make([]string, len(args)+1)
			newArgs[0] = "version"
			copy(newArgs[1:], args)
			args = newArgs
			break
		}
	}

	// Build the commands to include in the help now. This is pretty...
	// tedious, but we don't have a better way at the moment.
	commandsInclude := make([]string, 0, len(commands))
	for k := range commands {
		switch k {
		case "token-disk":
		default:
			commandsInclude = append(commandsInclude, k)
		}
	}

	cli := &cli.CLI{
		Args:         args,
		Commands:     commands,
		Name:         "gitx",
		Autocomplete: true,
		HelpFunc:     cli.FilteredHelpFunc(commandsInclude, HelpFunc),
	}

	exitCode, err := cli.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err.Error())
		return 1
	}

	return exitCode
}
