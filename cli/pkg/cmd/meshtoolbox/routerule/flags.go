package routerule

import (
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/spf13/cobra"
)

func AddBaseFlags(cmd *cobra.Command, opts *options.Options) {
	rrOpts := &(opts.Create).InputRoutingRule
	flags := cmd.Flags()

	flags.StringVar(&(rrOpts.TargetMesh).Name,
		"mesh",
		"",
		"The mesh that will be the target for this rule")

	flags.StringVarP(&(rrOpts.TargetMesh).Namespace,
		"namespace", "n",
		"",
		"The namespace for this routing rule. Defaults to \"default\"")

	flags.StringVar(&rrOpts.Sources,
		"sources",
		"",
		"Sources for this rule. Each entry consists of an upstream namespace and and upstream name, separated by a colon.")

	flags.StringVar(&rrOpts.Destinations,
		"destinations",
		"",
		"Destinations for this rule. Same format as for 'sources'")

	flags.BoolVar(&rrOpts.OverrideExisting,
		"override",
		false,
		"If set to \"true\", the command will override any existing routing rule that matches the given namespace and name")

	rrOpts.Matchers = *flags.StringArrayP("matchers",
		"m",
		nil,
		"Matcher for this rule")
}

func AddTimeoutFlags(cmd *cobra.Command, opts *options.Options) {
	rrOpts := &(opts.Create).InputRoutingRule
	flags := cmd.Flags()
	flags.StringVar(&(rrOpts.Timeout).Seconds,
		"route.timeout.seconds",
		"",
		"timeout time in seconds")
	flags.StringVar(&(rrOpts.Timeout).Nanos,
		"route.timeout.nanos",
		"",
		"timeout time in nanoseconds")
}

func AddRetryFlags(cmd *cobra.Command, opts *options.Options) {
	rrOpts := &(opts.Create).InputRoutingRule
	flags := cmd.Flags()
	flags.StringVar(&(rrOpts.Retry).Attempts,
		"route.retry.attempt",
		"",
		"number of times to retry")
	flags.StringVar(&(rrOpts.Retry.PerTryTimeout).Seconds,
		"route.retry.timeout.seconds",
		"",
		"retry timeout time in seconds")
	flags.StringVar(&(rrOpts.Retry.PerTryTimeout).Nanos,
		"route.retry.timeout.nanos",
		"",
		"retry timeout time in nanoseconds")
}

func AddFaultFlags(cmd *cobra.Command, opts *options.Options) {
	rrOpts := &(opts.Create).InputRoutingRule
	flags := cmd.Flags()
	// delay
	flags.StringVar(&(rrOpts.FaultInjection).DelayPercent,
		"fault.delay.percent",
		"",
		"Percentage of requests on which the delay will be injected (0-100).")
	flags.StringVar(&(rrOpts.FaultInjection).HttpDelayType,
		"fault.delay.type",
		"",
		"Type of delay (fixed or exponential).")
	flags.StringVar(&(rrOpts.FaultInjection).HttpDelayValue.Seconds,
		"fault.delay.value.seconds",
		"",
		"delay duration (seconds).")
	flags.StringVar(&(rrOpts.FaultInjection).HttpDelayValue.Nanos,
		"fault.delay.value.nanos",
		"",
		"delay duration (nanoseconds).")
	// abort
	flags.StringVar(&(rrOpts.FaultInjection).AbortPercent,
		"fault.abort.percent",
		"",
		"Percentage of requests on which the abort will be injected (0-100).")
	flags.StringVar(&(rrOpts.FaultInjection).ErrorType,
		"fault.abort.type",
		"",
		"Type of error (http, http2, or grpc).")
	flags.StringVar(&(rrOpts.FaultInjection).ErrorMessage,
		"fault.abort.message",
		"",
		"Error message (int for type=http errors, string otherwise).")
}
