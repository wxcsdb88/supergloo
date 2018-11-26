package get

import (
	"fmt"
	"strings"

	"github.com/solo-io/supergloo/cli/pkg/cmd/get/info"
	"github.com/solo-io/supergloo/cli/pkg/common"

	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

var supportedOutputFormats = []string{"wide", "yaml"}

func Cmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: `Display one or many supergloo resources`,
		Long:  `Display one or many supergloo resources`,
		Args:  cobra.RangeArgs(1, 2),
		Run: func(c *cobra.Command, args []string) {
			if err := get(args, opts); err != nil {
				fmt.Println(err)
			}
		},
	}
	getOpts := &opts.Get
	pFlags := cmd.Flags()
	pFlags.StringVarP(&getOpts.Output, "output", "o", "",
		"Output format. Options include: \n"+strings.Join(supportedOutputFormats, "|"))
	return cmd
}

func get(args []string, opts *options.Options) error {

	infoClient, err := info.NewClient()
	if err != nil {
		return err
	}

	if err := ensureParameters(infoClient, opts, args); err != nil {
		return err
	}

	return getResource(infoClient, opts.Get)
}

func ensureParameters(infoClient info.SuperglooInfoClient, opts *options.Options, args []string) error {
	gOpts := &opts.Get

	// Argument count is validated by cobra.RangeArgs
	// first arg is the resource type
	gOpts.Type = args[0]
	// second arg is the resource name (optional)
	gOpts.Name = ""
	if len(args) == 2 {
		gOpts.Name = args[1]
	}

	// Get available resource types
	resourceTypes, err := infoClient.ListResourceTypes()
	if err != nil {
		return err
	}

	// Validate input resource type
	if !common.Contains(resourceTypes, gOpts.Type) {
		return errors.Errorf(common.UnknownResourceTypeMsg, gOpts.Type)
	}

	// Output format is set by a flag
	if gOpts.Output != "" && !common.Contains(supportedOutputFormats, gOpts.Output) {
		return errors.Errorf(common.UnknownOutputFormat, gOpts.Output, strings.Join(supportedOutputFormats, "|"))
	}

	return nil
}

func getResource(infoClient info.SuperglooInfoClient, gOpts options.Get) error {

	// Fetch the resource information
	err := infoClient.ListResources(gOpts)
	if err != nil {
		return err
	}
	return nil
}
