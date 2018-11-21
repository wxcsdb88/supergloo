package cmd

import (
	"github.com/pkg/errors"
	"github.com/solo-io/supergloo/cli/pkg/cmd/create"
	"github.com/solo-io/supergloo/cli/pkg/cmd/get"
	"github.com/solo-io/supergloo/cli/pkg/cmd/ingresstoolbox"
	"github.com/solo-io/supergloo/cli/pkg/cmd/install"
	"github.com/solo-io/supergloo/cli/pkg/cmd/meshtoolbox"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

var opts options.Options

func App(version string) *cobra.Command {
	app := &cobra.Command{
		Use:   "supergloo",
		Short: "manage mesh resources with supergloo",
		Long: `supergloo configures resources used by Supergloo server.
	Find more information at https://supergloo.solo.io`,
		Version: version,
	}

	pflags := app.PersistentFlags()
	pflags.BoolVarP(&opts.Top.Static, "static", "s", false, "disable interactive mode")

	app.SuggestionsMinimumDistance = 1
	app.AddCommand(
		install.Cmd(&opts),

		get.Cmd(&opts),
		create.Cmd(&opts),
		meshtoolbox.FaultInjection(&opts),
		meshtoolbox.LoadBalancing(&opts),
		meshtoolbox.Retries(&opts),
		ingresstoolbox.FortifyIngress(&opts),
		ingresstoolbox.AddRoute(&opts),
	)

	return app
}

func withInit(command *cobra.Command, initFunc func() error) {
	if command.RunE != nil {
		return
	}
	original := command.RunE
	command.RunE = func(cmd *cobra.Command, args []string) error {
		if err := initFunc(); err != nil {
			return errors.Wrapf(err, "initialization error")
		}
		return original(cmd, args)
	}
}
