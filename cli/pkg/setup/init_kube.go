package setup

import (
	"github.com/pkg/errors"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
)

func InitKubeOptions(opts *options.Options) error {
	if err := InitCache(opts); err != nil {
		return errors.Wrapf(err, "setting up cache")
	}

	if err := Init(opts); err != nil {
		return errors.Wrap(err, "Error during initialization.")
	}
	return nil
}
