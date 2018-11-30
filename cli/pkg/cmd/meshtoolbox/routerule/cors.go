package routerule

import (
	"fmt"
	"strings"

	types "github.com/gogo/protobuf/types"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/solo-io/supergloo/cli/pkg/common/iutil"
	v1alpha3 "github.com/solo-io/supergloo/pkg/api/external/istio/networking/v1alpha3"
)

func EnsureCors(irOpts *options.InputCors, opts *options.Options) error {
	cOpts := &(opts.Create.InputRoutingRule).Cors
	// initialize the field
	target := &v1alpha3.CorsPolicy{}

	if err := ensureCsv("Please specify the allowed origins (comma-separated list)", cOpts.AllowOrigin, &target.AllowOrigin, opts.Top.Static); err != nil {
		return err
	}
	if err := ensureCsv("Please specify the allowed methods (comma-separated list)", cOpts.AllowMethods, &target.AllowMethods, opts.Top.Static); err != nil {
		return err
	}
	if err := ensureCsv("Please specify the allowed headers (comma-separated list)", cOpts.AllowHeaders, &target.AllowHeaders, opts.Top.Static); err != nil {
		return err
	}
	if err := ensureCsv("Please specify the exposed headers (comma-separated list)", cOpts.ExposeHeaders, &target.ExposeHeaders, opts.Top.Static); err != nil {
		return err
	}
	target.MaxAge = &types.Duration{}
	if err := EnsureDuration("Please specify the max age", &cOpts.MaxAge, target.MaxAge, opts); err != nil {
		return err
	}

	opts.MeshTool.RoutingRule.CorsPolicy = target
	return nil
}

func ensureCsv(message string, source string, target *[]string, staticMode bool) error {
	if staticMode && source == "" {
		return fmt.Errorf(message)
	}
	if !staticMode {
		if err := iutil.GetStringInput(message, &source); err != nil {
			return err
		}
	}
	parts := strings.Split(source, ",")
	*target = parts
	return nil
}
