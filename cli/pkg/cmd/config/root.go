package config

import (
	"github.com/solo-io/supergloo/cli/pkg/cmd/config/ca"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func Cmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: `Configure mesh resources`,
		Long:  `Configure mesh resources`,
		Args:  cobra.ExactArgs(1), // TODO: for now allow only stdin creation, no file
		Run: func(c *cobra.Command, args []string) {
		},
	}

	cmd.AddCommand(
		ca.Cmd(opts),
	)

	return cmd
}
