package translator

import (
	"context"

	"github.com/solo-io/supergloo/pkg/api/v1"
)

type Syncer struct{}

func (sync *Syncer) Sync(ctx context.Context, snap *v1.TranslatorSnapshot) error {
	panic("implement me")
}
