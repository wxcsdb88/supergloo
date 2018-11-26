package istio

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/solo-io/supergloo/pkg/api/v1"
	"k8s.io/client-go/kubernetes"

	security "github.com/openshift/client-go/security/clientset/versioned"
	kubemeta "k8s.io/apimachinery/pkg/apis/meta/v1"

	crdClient "k8s.io/apiextensions-apiserver/pkg/client/clientset/internalclientset/typed/apiextensions/internalversion"
)

const (
	CrbName          = "istio-crb"
	defaultNamespace = "istio-system"
)

type IstioInstaller struct {
	SecurityClient *security.Clientset
	CrdClient      *crdClient.ApiextensionsClient
}

func (c *IstioInstaller) GetDefaultNamespace() string {
	return defaultNamespace
}

func (c *IstioInstaller) GetCrbName() string {
	return CrbName
}

func (c *IstioInstaller) GetOverridesYaml(install *v1.Install) string {
	return getOverrides(install.Encryption)
}

func getOverrides(encryption *v1.Encryption) string {
	selfSigned := false
	mtlsEnabled := false
	if encryption != nil {
		if encryption.TlsEnabled {
			mtlsEnabled = true
			if encryption.Secret != nil {
				selfSigned = true
			}
		}
	}
	selfSignedString := strconv.FormatBool(selfSigned)
	tlsEnabledString := strconv.FormatBool(mtlsEnabled)
	overridesWithMtlsFlag := strings.Replace(overridesYaml, "@@MTLS_ENABLED@@", tlsEnabledString, -1)
	return strings.Replace(overridesWithMtlsFlag, "@@SELF_SIGNED@@", selfSignedString, -1)
}

var overridesYaml = `
global.mtls.enabled: @@MTLS_ENABLED@@
security.selfSigned: @@SELF_SIGNED@@
`

func (c *IstioInstaller) DoPostHelmInstall(install *v1.Install, kube *kubernetes.Clientset, releaseName string) error {
	return nil
}

func (c *IstioInstaller) DoPreHelmInstall() error {
	if c.SecurityClient == nil {
		return nil
	}
	return c.AddSccToUsers(
		"default",
		"istio-ingress-service-account",
		"prometheus",
		"istio-egressgateway-service-account",
		"istio-citadel-service-account",
		"istio-ingressgateway-service-account",
		"istio-cleanup-old-ca-service-account",
		"istio-mixer-post-install-account",
		"istio-mixer-service-account",
		"istio-pilot-service-account",
		"istio-sidecar-injector-service-account",
		"istio-galley-service-account")
}

// TODO: something like this should enable minishift installs to succeed, but this isn't right. The correct steps are
//       to run "oc adm policy add-scc-to-user anyuid -z %s -n istio-system" for each of the user accounts above
//       maybe the issue is not specifying the namespace?
func (c *IstioInstaller) AddSccToUsers(users ...string) error {
	anyuid, err := c.SecurityClient.SecurityV1().SecurityContextConstraints().Get("anyuid", kubemeta.GetOptions{})
	if err != nil {
		return err
	}
	newUsers := append(anyuid.Users, users...)
	anyuid.Users = newUsers
	_, err = c.SecurityClient.SecurityV1().SecurityContextConstraints().Update(anyuid)
	return err
}

func (c *IstioInstaller) DoPostHelmUninstall() error {
	// TODO: this will break if there are more than one installs using these CRDs
	if err := c.deleteIstioCrds(); err != nil {
		return err
	}
	return nil
}

func (c *IstioInstaller) deleteIstioCrds() error {
	if c.CrdClient == nil {
		return errors.Errorf("Crd client not provided")
	}
	crdList, err := c.CrdClient.CustomResourceDefinitions().List(kubemeta.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "Error getting crds")
	}
	for _, crd := range crdList.Items {
		//TODO: use labels
		if strings.Contains(crd.Name, "istio.io") {
			err = c.CrdClient.CustomResourceDefinitions().Delete(crd.Name, &kubemeta.DeleteOptions{})
			if err != nil {
				return errors.Wrap(err, "Error deleting crd")
			}
		}
	}
	return nil
}
