package translator_test

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/supergloo/pkg/api/v1"

	. "github.com/solo-io/supergloo/pkg/translator"
)

var _ = Describe("Syncer", func() {
	It("works", func() {
		s := &Syncer{}
		err := s.Sync(context.TODO(), &v1.TranslatorSnapshot{
			Meshes: v1.MeshesByNamespace{
				"mymesh": {
					{
						Metadata: core.Metadata{Namespace:"mymesh", Name:"yyaa"},
					},
				},
			},
		})
		Expect(err).NotTo(HaveOccurred())
	})
})
