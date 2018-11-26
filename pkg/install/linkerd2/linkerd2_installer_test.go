package linkerd2_test

import (
	"context"
	"os"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/supergloo/pkg/install"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/supergloo/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/supergloo/pkg/api/v1"
)

/*
Smoke test for installing and uninstalling linkerd2
*/
var _ = Describe("Linkerd2 Installer", func() {

	superglooNamespace := "supergloo-system" // this needs to be made before running tests
	meshName := "test-linkerd-mesh"

	path := os.Getenv("HELM_CHART_PATH")
	if path == "" {
		path = "https://s3.amazonaws.com/supergloo.solo.io/linkerd2-0.1.0.tgz"
	}

	getSnapshot := func(install bool) *v1.InstallSnapshot {
		return &v1.InstallSnapshot{
			Installs: v1.InstallsByNamespace{
				superglooNamespace: v1.InstallList{
					&v1.Install{
						Metadata: core.Metadata{
							Namespace: superglooNamespace,
							Name:      meshName,
						},
						MeshType: &v1.Install_Linkerd2{
							Linkerd2: &v1.Linkerd2{
								// not specifying install namespace since it is hard coded in chart
							},
						},
						ChartLocator: &v1.HelmChartLocator{
							Kind: &v1.HelmChartLocator_ChartPath{
								ChartPath: &v1.HelmChartPath{
									Path: path,
								},
							},
						},
						Enabled: &types.BoolValue{
							Value: install,
						},
					},
				},
			},
		}
	}

	kubeCache := kube.NewKubeCache()

	var meshClient v1.MeshClient
	var syncer install.InstallSyncer

	BeforeEach(func() {
		util.TryCreateNamespace("supergloo-system")
		meshClient = util.GetMeshClient(kubeCache)
		syncer = install.InstallSyncer{
			Kube:       util.GetKubeClient(),
			MeshClient: meshClient,
		}
	})

	AfterEach(func() {
		util.TerminateNamespaceBlocking("supergloo-system")

		// just in case
		meshClient.Delete(superglooNamespace, meshName, clients.DeleteOpts{})
		util.UninstallHelmRelease(meshName)
		util.TerminateNamespaceBlocking("linkerd") //hard-coded in chart
	})

	It("Can install and uninstall linkerd2", func() {
		snap := getSnapshot(true)
		err := syncer.Sync(context.TODO(), snap)
		Expect(err).NotTo(HaveOccurred())
		util.WaitForAvailablePods("linkerd")

		snap = getSnapshot(false)
		err = syncer.Sync(context.TODO(), snap)
		Expect(err).NotTo(HaveOccurred())

		// Validate everything cleaned up
		util.WaitForTerminatedNamespace("linkerd")
		Expect(util.HelmReleaseDoesntExist(meshName)).To(BeTrue())

		mesh, err := meshClient.Read(superglooNamespace, meshName, clients.ReadOpts{})
		Expect(mesh).To(BeNil())
		Expect(err).ToNot(BeNil())

	})
})
