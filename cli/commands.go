package cli

import (
	"os"

	"github.com/go-almeria/gitter/command"
	"github.com/go-almeria/gitter/meta"
	"github.com/go-almeria/gitter/version"

	"github.com/mitchellh/cli"
)

// Commands returns the mapping of CLI commands for Gitter. The meta
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
