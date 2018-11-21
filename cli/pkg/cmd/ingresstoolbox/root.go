package ingresstoolbox

import (
	"fmt"
	"github.com/solo-io/supergloo/cli/pkg/setup"

	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func FortifyIngress(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fortify-ingress",
		Short: `Configure ingress security parameters`,
		Long:  `Configure ingress security parameters`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := setup.InitKubeOptions(opts); err != nil {
				return err
			}
			return ingressToolPlaceholder(opts)
		},
	}
	linkIngressToolFlags(cmd, opts)
	return cmd
}

func AddRoute(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-route",
		Short: `Define new route`,
		Long:  `Define new route`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := setup.InitKubeOptions(opts); err != nil {
				return err
			}
			return ingressToolPlaceholder(opts)
		},
	}
	linkIngressToolFlags(cmd, opts)
	return cmd
}

func ingressToolPlaceholder(opts *options.Options) error {
	fmt.Println("not implemented")
	return nil
}

func linkIngressToolFlags(cmd *cobra.Command, opts *options.Options) {
	pflags := cmd.PersistentFlags()
	pflags.StringVar(&opts.IngressTool.IngressId, "ingressid", "", "ingress to modify")
	pflags.StringVar(&opts.IngressTool.RouteId, "routeid", "", "route to modify")
}
