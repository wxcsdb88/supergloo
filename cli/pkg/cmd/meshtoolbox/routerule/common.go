package routerule

import (
	"fmt"
	"strconv"

	types "github.com/gogo/protobuf/types"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/solo-io/supergloo/cli/pkg/common/iutil"
)

func EnsureDuration(rootMessage string, durOpts *options.InputDuration, targetDur *types.Duration, opts *options.Options) error {
	dur := types.Duration{}
	if !opts.Top.Static {
		err := iutil.GetStringInput(fmt.Sprintf("%v (seconds)", rootMessage), &durOpts.Seconds)
		if err != nil {
			return err
		}
		err = iutil.GetStringInput(fmt.Sprintf("%v (nanoseconds)", rootMessage), &durOpts.Nanos)
		if err != nil {
			return err
		}
	}
	// if not in interactive mode, timeout values will have already been passed
	if durOpts.Seconds != "" {
		sec, err := strconv.Atoi(durOpts.Seconds)
		if err != nil {
			return err
		}
		dur.Seconds = int64(sec)
	}
	if durOpts.Nanos != "" {
		nanos, err := strconv.Atoi(durOpts.Nanos)
		if err != nil {
			return err
		}
		dur.Nanos = int32(nanos)
	}
	*targetDur = dur
	return nil
}

// EnsurePercentage transforms a source string to a target int
// If not present, it promts the user for input with the given message
// Errors on invalid input
func EnsurePercentage(message string, source *string, target *int32, opts *options.Options) error {
	if !opts.Top.Static {
		if err := iutil.GetStringInput(message, source); err != nil {
			return err
		}
	}
	if *source != "" {
		percentage, err := strconv.Atoi(*source)
		if err != nil {
			return err
		}
		*target = int32(percentage)
	}
	return nil
}
