package meshtoolbox

import (
	"fmt"

	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func FaultInjection(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fault-injection",
		Short: `Stress test your mesh with faults`,
		Long:  `Stress test your mesh with faults`,
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("this feature will be available in 2019")
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
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("this feature will be available in 2019")
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
		Run: func(c *cobra.Command, args []string) {
			meshToolPlaceholder(opts)
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
