package cli

import (
	"os"

	"github.com/go-almeria/gitx/command"
	"github.com/go-almeria/gitx/meta"
	"github.com/go-almeria/gitx/version"

	"github.com/mitchellh/cli"
)

// Commands returns the mapping of CLI commands for Gitx. The meta
// parameter lets you set meta options for all commands.
func Commands(metaPtr *meta.Meta) map[string]cli.CommandFactory {

	if metaPtr == nil {
		metaPtr = &meta.Meta{}
	}

	if metaPtr.Ui == nil {
		metaPtr.Ui = &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		}
	}

	return map[string]cli.CommandFactory{
		"version": func() (cli.Command, error) {
			versionInfo := version.GetVersion()

			return &command.VersionCommand{
				VersionInfo: versionInfo,
				Ui:          metaPtr.Ui,
			}, nil
		},
	}
}
