package meshtoolbox

import (
	"fmt"

	"github.com/solo-io/supergloo/cli/pkg/cmd/meshtoolbox/mtls"
	"github.com/solo-io/supergloo/cli/pkg/cmd/meshtoolbox/policy"
	"github.com/solo-io/supergloo/cli/pkg/cmd/meshtoolbox/routerule"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func FaultInjection(opts *options.Options) *cobra.Command {
	rrOpts := &(opts.Create).InputRoutingRule
	cmd := &cobra.Command{
		Use:   "fault-injection",
		Short: `Stress test your mesh with faults`,
		Long:  `Stress test your mesh with faults`,
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rrOpts.RouteName = args[0]
			if err := routerule.CreateRoutingRule(routerule.FaultInjection_Rule, opts); err != nil {
				return err
			}
			fmt.Printf("Created fault injection routing rule [%v] in namespace [%v]\n", args[0], rrOpts.TargetMesh.Namespace)
			return nil
		},
	}
	linkMeshToolFlags(cmd, opts)
	routerule.AddFaultFlags(cmd, opts)
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

func Policy(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: `Apply a policy`,
		Long:  `Apply, update, or remove a policy`,
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
		},
	}
	linkMeshToolFlags(cmd, opts)
	policy.LinkPolicyFlags(cmd, opts)
	cmd.AddCommand(
		policy.Add(opts),
		policy.Remove(opts),
		policy.Clear(opts),
	)
	return cmd
}

func ToggleMtls(opts *options.Options) *cobra.Command {
	cmd := mtls.Root(opts)
	linkMeshToolFlags(cmd, opts)
	return cmd
}

func linkMeshToolFlags(cmd *cobra.Command, opts *options.Options) {
	meshRef := &(opts.MeshTool).Mesh
	pflags := cmd.PersistentFlags()
	pflags.StringVar(&meshRef.Name, "mesh.name", "", "name of mesh to update")
	pflags.StringVar(&meshRef.Namespace, "mesh.namespace", "", "namespace of mesh to update")
	pflags.StringVar(&opts.MeshTool.ServiceId, "serviceid", "", "service to modify")
}

func meshToolPlaceholder(opts *options.Options) {
	fmt.Println("this mesh feature will be available in 2019")
}
