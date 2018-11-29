package create

import (
	"fmt"

	"github.com/solo-io/supergloo/cli/pkg/common"

	types "github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/supergloo/cli/pkg/cmd/meshtoolbox/routerule"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/solo-io/supergloo/cli/pkg/nsutil"
	glooV1 "github.com/solo-io/supergloo/pkg/api/external/gloo/v1"
	superglooV1 "github.com/solo-io/supergloo/pkg/api/v1"
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
			if err := createRoutingRule(args[0], opts); err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Created routing rule [%v] in namespace [%v]\n", args[0], rrOpts.TargetMesh.Namespace)
		},
	}

	routerule.AddBaseFlags(cmd, opts)

	return cmd
}

func createRoutingRule(routeName string, opts *options.Options) error {
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
	// Initialize the root of our RoutingRule with the minimal required params
	// TODO(mitchdraft) move these fields out s.t. they are populated by the ensure methods
	opts.MeshTool.RoutingRule = superglooV1.RoutingRule{
		Metadata: core.Metadata{
			Name:      routeName,
			Namespace: rrOpts.TargetMesh.Namespace,
		},
		TargetMesh:      &rrOpts.TargetMesh,
		Sources:         opts.MeshTool.RoutingRule.Sources,
		Destinations:    opts.MeshTool.RoutingRule.Destinations,
		RequestMatchers: opts.MeshTool.RoutingRule.RequestMatchers,
	}

	// TODO(mitchdraft) gate this behind setting (so that it can be called from a top-level command)
	opts.MeshTool.RoutingRule.Timeout = &types.Duration{}
	if err := routerule.EnsureTimeout(&(opts.Create.InputRoutingRule).Timeout, opts.MeshTool.RoutingRule.Timeout, opts); err != nil {
		return err
	}
	if err := routerule.EnsureRetry(&(opts.Create.InputRoutingRule).Retry, opts); err != nil {
		return err
	}

	rrClient, err := common.GetRoutingRuleClient()
	if err != nil {
		return err
	}
	_, err = (*rrClient).Write(&(opts.MeshTool).RoutingRule, clients.WriteOpts{OverwriteExisting: rrOpts.OverrideExisting})
	return err
}
