package create

import (
	"fmt"

	"github.com/solo-io/supergloo/cli/pkg/common"

	// types "github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	// "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/supergloo/cli/pkg/cmd/meshtoolbox/routerule"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/solo-io/supergloo/cli/pkg/nsutil"
	glooV1 "github.com/solo-io/supergloo/pkg/api/external/gloo/v1"

	// superglooV1 "github.com/solo-io/supergloo/pkg/api/v1"
	"github.com/spf13/cobra"
)

func RoutingRuleCmd(opts *options.Options) *cobra.Command {
	rrOpts := &(opts.Create).InputRoutingRule
	cmd := &cobra.Command{
		Use:   "routingrule",
		Short: `Create a route rule with the given name`,
		Long:  `Create a route rule with the given name`,
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			rrOpts.RouteName = args[0]
			if err := createRoutingRule(opts); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Created routing rule [%v] in namespace [%v]\n", args[0], rrOpts.TargetMesh.Namespace)
		},
	}

	routerule.AddBaseFlags(cmd, opts)
	routerule.AddTimeoutFlags(cmd, opts)
	routerule.AddRetryFlags(cmd, opts)
	routerule.AddFaultFlags(cmd, opts)

	return cmd
}

func createRoutingRule(opts *options.Options) error {
	rrOpts := &(opts.Create).InputRoutingRule

	// all operations require a target mesh spec
	if err := nsutil.EnsureMesh(&rrOpts.TargetMesh, opts); err != nil {
		return err
	}

	// Validate source and destination upstreams
	if err := routerule.EnsureUpstreams(opts); err != nil {
		return err
	}

	// Validate matchers
	opts.MeshTool.RoutingRule.RequestMatchers = []*glooV1.Matcher{}
	if rrOpts.Matchers != nil {
		if err := routerule.ValidateMatchers(rrOpts.Matchers, opts.MeshTool.RoutingRule.RequestMatchers); err != nil {
			return err
		}

	}

	opts.Create.InputRoutingRule.ActiveTypes = routerule.GenerateActiveRuleList(routerule.USE_ALL_ROUTING_RULES)
	fmt.Println("opts.Create.InputRoutingRule.ActiveTypes")
	fmt.Println(opts.Create.InputRoutingRule.ActiveTypes)
	if err := routerule.EnsureActiveRoutingRuleTypes(&(opts.Create.InputRoutingRule).ActiveTypes, opts.Top.Static); err != nil {
		return err
	}
	if err := routerule.AssembleRoutingRule(opts.Create.InputRoutingRule.ActiveTypes, opts); err != nil {
		return err
	}

	rrClient, err := common.GetRoutingRuleClient()
	if err != nil {
		return err
	}
	_, err = (*rrClient).Write(&(opts.MeshTool).RoutingRule, clients.WriteOpts{OverwriteExisting: rrOpts.OverrideExisting})
	return err
}
