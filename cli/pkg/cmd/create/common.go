package create

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/supergloo/cli/pkg/constants"
	glooV1 "github.com/solo-io/supergloo/pkg/api/external/gloo/v1"
	superglooV1 "github.com/solo-io/supergloo/pkg/api/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getUpstreamClient() (*glooV1.UpstreamClient, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}

	upstreamClient, err := glooV1.NewUpstreamClient(&factory.KubeResourceClientFactory{
		Crd:         glooV1.UpstreamCrd,
		Cfg:         config,
		SharedCache: kube.NewKubeCache(),
	})
	if err != nil {
		return nil, err
	}
	if err = upstreamClient.Register(); err != nil {
		return nil, err
	}
	return &upstreamClient, nil
}

func getMeshClient() (*superglooV1.MeshClient, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}
	meshClient, err := superglooV1.NewMeshClient(&factory.KubeResourceClientFactory{
		Crd:         superglooV1.MeshCrd,
		Cfg:         config,
		SharedCache: kube.NewKubeCache(),
	})
	if err != nil {
		return nil, err
	}
	if err = meshClient.Register(); err != nil {
		return nil, err
	}
	return &meshClient, nil
}

func getRoutingRuleClient() (*superglooV1.RoutingRuleClient, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}
	rrClient, err := superglooV1.NewRoutingRuleClient(&factory.KubeResourceClientFactory{
		Crd:         superglooV1.RoutingRuleCrd,
		Cfg:         config,
		SharedCache: kube.NewKubeCache(),
	})
	if err != nil {
		return nil, err
	}
	if err = rrClient.Register(); err != nil {
		return nil, err
	}
	return &rrClient, nil
}

func getKubernetesClient() (*kubernetes.Clientset, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return kubeClient, nil
}

func getKubernetesConfig() (*rest.Config, error) {
	config, err := kubeutils.GetConfig("", "")
	if err != nil {
		return nil, fmt.Errorf(constants.KubeConfigError, err)
	}
	return config, nil
}
