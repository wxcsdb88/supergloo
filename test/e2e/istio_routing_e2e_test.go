package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/test/helpers"
	testsetup "github.com/solo-io/solo-kit/test/setup"
	"github.com/solo-io/supergloo/pkg/api/v1"
	"github.com/solo-io/supergloo/pkg/setup"
	"os"
	"os/exec"
	"time"
)

var _ = Describe("istio routing E2e", func() {
	var namespace string
	path := os.Getenv("HELM_CHART_PATH")
	if path == "" {
		Skip("Set environment variable HELM_CHART_PATH")
	}
	BeforeEach(func() {
		namespace = helpers.RandString(8)
		err := testsetup.SetupKubeForTest(namespace)
		Expect(err).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		exec.Command("helm", "delete", "my-istio-install").Run()
		testsetup.TeardownKube(namespace)
	})

	It("works", func() {
		go setup.Main(namespace)

		// start discovery
		cmd := exec.Command(PathToUds, "--namespace", namespace)
		cmd.Env = os.Environ()
		_, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())

		meshes, routingRules, installClient, err := run()
		Expect(err).NotTo(HaveOccurred())

		installName := setupInstall(installClient, namespace, path)

		var ref *core.ResourceRef
		Eventually(func() (*core.ResourceRef, error) {
			mesh, err := meshes.Read(namespace, installName, clients.ReadOpts{})
			if err != nil {
				return nil, err
			}
			r := mesh.Metadata.Ref()
			ref = &r
			return ref, nil
		}, time.Minute * 2).Should(Not(BeNil()))

		setupRoutingRule(routingRules, namespace, ref)

		// ignored
		if false {
			meshMeta := core.Metadata{Name: "my-istio", Namespace: namespace}
			meshes.Delete(meshMeta.Namespace, meshMeta.Name, clients.DeleteOpts{})
			m1, err := meshes.Write(&v1.Mesh{
				Metadata: meshMeta,
				MeshType: &v1.Mesh_Istio{
					Istio: &v1.Istio{
						WatchNamespaces: []string{namespace},
					},
				},
				Encryption: &v1.Encryption{
					TlsEnabled: true,
				},
				//Encryption: &v1.Encryption{TlsEnabled: true},
			}, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			Expect(m1).NotTo(BeNil())
		}

	})
})

func setupInstall(installClient v1.InstallClient, namespace string, chartPath string) string {
	installMeta := core.Metadata{Name: "my-istio-install", Namespace: namespace}
	installClient.Delete(installMeta.Namespace, installMeta.Name, clients.DeleteOpts{})

	install1, err := installClient.Write(&v1.Install{
		Metadata: installMeta,
		MeshType: &v1.Install_Istio{
			Istio: &v1.Istio{
				InstallationNamespace: namespace,
			},
		},
		ChartLocator: &v1.HelmChartLocator{
			Kind: &v1.HelmChartLocator_ChartPath{
				ChartPath: &v1.HelmChartPath{
					Path: chartPath,
				},
			},
		},
		Encryption: &v1.Encryption{
			TlsEnabled: true,
		},
	}, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(install1).NotTo(BeNil())
	return install1.Metadata.Name
}

func setupRoutingRule(routingRules v1.RoutingRuleClient, namespace string, targetMesh *core.ResourceRef) {
	rrMeta := core.Metadata{Name: "my-istio-rr", Namespace: namespace}
	routingRules.Delete(rrMeta.Namespace, rrMeta.Name, clients.DeleteOpts{})
	rr1, err := routingRules.Write(&v1.RoutingRule{
		Metadata:   rrMeta,
		TargetMesh: targetMesh,
		Destinations: []*core.ResourceRef{{
			Name:      "default-reviews-9080",
			Namespace: "gloo-system",
		}},
		TrafficShifting: &v1.TrafficShifting{
			Destinations: []*v1.WeightedDestination{
				{
					Upstream: &core.ResourceRef{
						Name:      "default-reviews-v1-9080",
						Namespace: "gloo-system",
					},
					Weight: 100,
				},
			},
		},
	}, clients.WriteOpts{})
	Expect(err).NotTo(HaveOccurred())
	Expect(rr1).NotTo(BeNil())
}

func run() (v1.MeshClient, v1.RoutingRuleClient, v1.InstallClient, error) {
	kubeCache := kube.NewKubeCache()
	restConfig, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, nil, nil, err
	}
	meshClient, err := v1.NewMeshClient(&factory.KubeResourceClientFactory{
		Crd:         v1.MeshCrd,
		Cfg:         restConfig,
		SharedCache: kubeCache,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	if err := meshClient.Register(); err != nil {
		return nil, nil, nil, err
	}

	routingRuleClient, err := v1.NewRoutingRuleClient(&factory.KubeResourceClientFactory{
		Crd:         v1.RoutingRuleCrd,
		Cfg:         restConfig,
		SharedCache: kubeCache,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	if err := routingRuleClient.Register(); err != nil {
		return nil, nil, nil, err
	}

	installClient, err := v1.NewInstallClient(&factory.KubeResourceClientFactory{
		Crd:         v1.InstallCrd,
		Cfg:         restConfig,
		SharedCache: kubeCache,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	if err := routingRuleClient.Register(); err != nil {
		return nil, nil, nil, err
	}
	return meshClient, routingRuleClient, installClient, nil
}
