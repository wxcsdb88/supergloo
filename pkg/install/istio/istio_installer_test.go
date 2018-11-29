package istio_test

import (
	"context"
	"os"

	"github.com/solo-io/solo-kit/test/helpers"

	"github.com/gogo/protobuf/types"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/supergloo/pkg/api/v1"
	"github.com/solo-io/supergloo/pkg/install"
	"github.com/solo-io/supergloo/pkg/install/istio"
	"github.com/solo-io/supergloo/test/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

/*
Smoke test for installing and uninstalling istio.
*/
var _ = Describe("Istio Installer", func() {

	superglooNamespace := "supergloo-system"

	path := os.Getenv("HELM_CHART_PATH")
	if path == "" {
		path = "https://s3.amazonaws.com/supergloo.solo.io/istio-1.0.3.tgz"
	}

	getSnapshot := func(install bool, installNamespace string, meshName string) *v1.InstallSnapshot {
		return &v1.InstallSnapshot{
			Installs: v1.InstallsByNamespace{
				superglooNamespace: v1.InstallList{
					&v1.Install{
						Metadata: core.Metadata{
							Namespace: superglooNamespace,
							Name:      meshName,
						},
						MeshType: &v1.Install_Istio{
							Istio: &v1.Istio{
								InstallationNamespace: installNamespace,
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
	var installNamespace string
	var meshName string

	BeforeEach(func() {
		randStr := helpers.RandString(8)
		installNamespace = "istio-install-test-" + randStr
		meshName = "istio-mesh-test-" + helpers.RandString(8)
		util.TryCreateNamespace(superglooNamespace)
		meshClient = util.GetMeshClient(kubeCache)
		syncer = install.InstallSyncer{
			Kube:       util.GetKubeClient(),
			MeshClient: meshClient,
			ApiExts:    util.GetApiExtsClient(),
		}
	})

	AfterEach(func() {
		// just in case
		meshClient.Delete(superglooNamespace, meshName, clients.DeleteOpts{})
		util.UninstallHelmRelease(meshName)
		util.TerminateNamespace(installNamespace) // random named, no need to block
		util.DeleteCrb(istio.CrbName)
		util.TerminateNamespaceBlocking(superglooNamespace)
	})

	It("Can install and uninstall istio", func() {
		snap := getSnapshot(true, installNamespace, meshName)
		err := syncer.Sync(context.TODO(), snap)
		Expect(err).NotTo(HaveOccurred())
		Expect(util.WaitForAvailablePods(installNamespace)).To(BeEquivalentTo(9))

		snap = getSnapshot(false, installNamespace, meshName)
		err = syncer.Sync(context.TODO(), snap)
		Expect(err).NotTo(HaveOccurred())

		// validate everything got cleaned up
		util.WaitForTerminatedNamespace(installNamespace)
		Expect(util.HelmReleaseDoesntExist(meshName)).To(BeTrue())
		Expect(util.CrbDoesntExist(istio.CrbName)).To(BeTrue())

		mesh, err := meshClient.Read(superglooNamespace, meshName, clients.ReadOpts{})
		Expect(mesh).To(BeNil())
		Expect(err).ToNot(BeNil())
	})
})
