package meshtoolbox

import (
	"fmt"
	"github.com/solo-io/supergloo/cli/pkg/setup"

	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func FaultInjection(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fault-injection",
		Short: `Stress test your mesh with faults`,
		Long:  `Stress test your mesh with faults`,
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Println("not implemented")
			return nil
		},
	}
	linkMeshToolFlags(cmd, opts)
	return cmd
}

func LoadBalancing(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "load-balancing",
		Short: `Specify traffic distribution`,
		Long:  `Specify traffic distribution`,
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Println("not implemented")
			return nil
		},
	}
	linkMeshToolFlags(cmd, opts)
	return cmd
}

func Retries(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retries",
		Short: `Configure retry parameters`,
		Long:  `Configure retry parameters`,
		RunE: func(c *cobra.Command, args []string) error {
			if err := setup.InitKubeOptions(opts); err != nil {
				return err
			}
			meshToolPlaceholder(opts)
			return nil
		},
	}
	linkMeshToolFlags(cmd, opts)
	return cmd
}

func linkMeshToolFlags(cmd *cobra.Command, opts *options.Options) {
	pflags := cmd.PersistentFlags()
	pflags.StringVar(&opts.MeshTool.MeshId, "meshid", "", "mesh to modify")
	pflags.StringVar(&opts.MeshTool.ServiceId, "serviceid", "", "service to modify")
}

func meshToolPlaceholder(opts *options.Options) {
	fmt.Println("this mesh feature will be available in 2019")
}
