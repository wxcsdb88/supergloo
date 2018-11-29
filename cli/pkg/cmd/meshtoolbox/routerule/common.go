package routerule

import (
	"strconv"

	types "github.com/gogo/protobuf/types"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/solo-io/supergloo/cli/pkg/common/iutil"
)

func EnsureTimeout(durOpts *options.InputDuration, targetDur *types.Duration, opts *options.Options) error {
	dur := types.Duration{}
	if !opts.Top.Static {
		err := iutil.GetStringInput("Please specify timeout duration (seconds)", &durOpts.Seconds)
		if err != nil {
			return err
		}
		err = iutil.GetStringInput("Please specify timeout duration (nanoseconds)", &durOpts.Nanos)
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
