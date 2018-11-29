package util

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/solo-io/supergloo/pkg/secret"

	"k8s.io/helm/pkg/proto/hapi/release"

	"github.com/solo-io/supergloo/pkg/install/helm"

	"github.com/hashicorp/consul/api"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/supergloo/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	gloo "github.com/solo-io/supergloo/pkg/api/external/gloo/v1"
	istiosecret "github.com/solo-io/supergloo/pkg/api/external/istio/encryption/v1"
	"k8s.io/client-go/kubernetes"

	kubecore "k8s.io/api/core/v1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kubemeta "k8s.io/apimachinery/pkg/apis/meta/v1"
	helmlib "k8s.io/helm/pkg/helm"
	helmkube "k8s.io/helm/pkg/kube"

	security "github.com/openshift/client-go/security/clientset/versioned"
	// love me google.
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	istiov1 "github.com/solo-io/supergloo/pkg/api/external/istio/encryption/v1"
)

var kubeConfig *rest.Config
var kubeClient *kubernetes.Clientset
var apiExtsClient apiexts.Interface
var upstreamClient gloo.UpstreamClient
var secretClient istiov1.IstioCacertsSecretClient

var testEcKey = "-----BEGIN PRIVATE KEY-----\nMIG2AgEAMBAGByqGSM49AgEGBSuBBAAiBIGeMIGbAgEBBDBoI1sMdiOTvBBdjWlS\nZ8qwNuK9xV4yKuboLZ4Sx/OBfy1eKZocxTKvnjLrHUe139uhZANiAAQMTIR56O8U\nTIqf6uUHM4i9mZYLj152up7elS06Gi6lk7IeUQDHxP0NnOnbhC7rmtOV6myLNApL\nQ92kZKg7qa8q7OY/4w1QfC4ch7zZKxjNkSIiuAx7V/lzF6FYDcqT3js=\n-----END PRIVATE KEY-----"
var TestEcRoot = "-----BEGIN CERTIFICATE-----\nMIIB7jCCAXUCCQC2t6Lqc2xnXDAKBggqhkjOPQQDAjBhMQswCQYDVQQGEwJVUzEW\nMBQGA1UECAwNTWFzc2FjaHVzZXR0czESMBAGA1UEBwwJQ2FtYnJpZGdlMQwwCgYD\nVQQKDANPcmcxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xODExMTgxMzQz\nMDJaFw0xOTExMTgxMzQzMDJaMGExCzAJBgNVBAYTAlVTMRYwFAYDVQQIDA1NYXNz\nYWNodXNldHRzMRIwEAYDVQQHDAlDYW1icmlkZ2UxDDAKBgNVBAoMA09yZzEYMBYG\nA1UEAwwPd3d3LmV4YW1wbGUuY29tMHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEDEyE\neejvFEyKn+rlBzOIvZmWC49edrqe3pUtOhoupZOyHlEAx8T9DZzp24Qu65rTleps\nizQKS0PdpGSoO6mvKuzmP+MNUHwuHIe82SsYzZEiIrgMe1f5cxehWA3Kk947MAoG\nCCqGSM49BAMCA2cAMGQCMCytVFc8sBdbM7DaBCz0N2ptdb0T7LFFfxDTzn4gjiDq\nVCd/3dct21TUWsthKXF2VgIwXEMI5EQiJ5kjR/y1KNBC9b4wfDiKRvG33jYe9gn6\ntzXUS00SoqG9D27/7aK71/xv\n-----END CERTIFICATE-----"
var testCertChain = ""

var testRsaCaKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAyzCxr/xu0zy5rVBiso9ffgl00bRKvB/HF4AX9/ytmZ6Hqsy1
3XIQk8/u/By9iCvVwXIMvyT0CbiJq/aPEj5mJUy0lzbrUs13oneXqrPXf7ir3Hzd
Rw+SBhXlsh9zAPZJXcF93DJU3GabPKwBvGJ0IVMJPIFCuDIPwW4kFAI7R/8A5LSd
PrFx6EyMXl7KM8jekC0y9DnTj83/fY72WcWX7YTpgZeBHAeeQOPTZ2KYbFal2gLs
ar69PgFS0TomESO9M14Yit7mzB1WDK2z9g3r+zLxENdJ5JG/ZskKe+TO4Diqi5OJ
t/h8yspS1ck8LJtCole9919umByg5oruflqIlQIDAQABAoIBAGZI8fnUinmd5R6B
C941XG3XFs6GAuUm3hNPcUFuGnntmv/5I0gBpqSyFO0nDqYg4u8Jma8TTCIkmnFN
ogIeFU+LiJFinR3GvwWzTE8rTz1FWoaY+M9P4ENd/I4pVLxUPuSKhfA2ChAVOupU
8F7D9Q/dfBXQQCT3VoUaC+FiqjL4HvIhji1zIqaqpK7fChGPraC/4WHwLMNzI0Zg
oDdAanwVygettvm6KD7AeKzhK94gX1PcnsOi3KuzQYvkenQE1M6/K7YtEc5qXCYf
QETj0UCzB55btgdF36BGoZXf0LwHqxys9ubfHuhwKBpY0xg2z4/4RXZNhfIDih3w
J3mihcECgYEA6FtQ0cfh0Zm03OPDpBGc6sdKxTw6aBDtE3KztfI2hl26xHQoeFqp
FmV/TbnExnppw+gWJtwx7IfvowUD8uRR2P0M2wGctWrMpnaEYTiLAPhXsj69HSM/
CYrh54KM0YWyjwNhtUzwbOTrh1jWtT9HV5e7ay9Atk3UWljuR74CFMUCgYEA392e
DVoDLE0XtbysmdlfSffhiQLP9sT8+bf/zYnr8Eq/4LWQoOtjEARbuCj3Oq7bP8IE
Vz45gT1mEE3IacC9neGwuEa6icBiuQi86NW8ilY/ZbOWrRPLOhk3zLiZ+yqkt+sN
cqWx0JkIh7IMKWI4dVQgk4I0jcFP7vNG/So4AZECgYEA426eSPgxHQwqcBuwn6Nt
yJCRq0UsljgbFfIr3Wfb3uFXsntQMZ3r67QlS1sONIgVhmBhbmARrcfQ0+xQ1SqO
wqnOL4AAd8K11iojoVXLGYP7ssieKysYxKpgPE8Yru0CveE9fkx0+OGJeM2IO5hY
qHAoTt3NpaPAuz5Y3XgqaVECgYA0TONS/TeGjxA9/jFY1Cbl8gp35vdNEKKFeM5D
Z7h+cAg56FE8tyFyqYIAGVoBFL7WO26mLzxiDEUfA/0Rb90c2JBfzO5hpleqIPd5
cg3VR+cRzI4kK16sWR3nLy2SN1k6OqjuovVS5Z3PjfI3bOIBz0C5FY9Pmt0g1yc7
mDRzcQKBgQCXWCZStbdjewaLd5u5Hhbw8tIWImMVfcfs3H1FN669LLpbARM8RtAa
8dYwDVHmWmevb/WX03LiSE+GCjCBO79fa1qc5RKAalqH/1OYxTuvYOeTUebSrg8+
lQFlP2OC4GGolKrN6HVWdxtf+F+SdjwX6qGCfYkXJRLYXIFSFjFeuw==
-----END RSA PRIVATE KEY-----`

var testRsaCaCert = `-----BEGIN CERTIFICATE-----
MIIDnzCCAoegAwIBAgIJAON1ifrBZ2/BMA0GCSqGSIb3DQEBCwUAMIGLMQswCQYD
VQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJU3Vubnl2YWxl
MQ4wDAYDVQQKDAVJc3RpbzENMAsGA1UECwwEVGVzdDEQMA4GA1UEAwwHUm9vdCBD
QTEiMCAGCSqGSIb3DQEJARYTdGVzdHJvb3RjYUBpc3Rpby5pbzAgFw0xODAxMjQx
OTE1NTFaGA8yMTE3MTIzMTE5MTU1MVowWTELMAkGA1UEBhMCVVMxEzARBgNVBAgT
CkNhbGlmb3JuaWExEjAQBgNVBAcTCVN1bm55dmFsZTEOMAwGA1UEChMFSXN0aW8x
ETAPBgNVBAMTCElzdGlvIENBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAyzCxr/xu0zy5rVBiso9ffgl00bRKvB/HF4AX9/ytmZ6Hqsy13XIQk8/u/By9
iCvVwXIMvyT0CbiJq/aPEj5mJUy0lzbrUs13oneXqrPXf7ir3HzdRw+SBhXlsh9z
APZJXcF93DJU3GabPKwBvGJ0IVMJPIFCuDIPwW4kFAI7R/8A5LSdPrFx6EyMXl7K
M8jekC0y9DnTj83/fY72WcWX7YTpgZeBHAeeQOPTZ2KYbFal2gLsar69PgFS0Tom
ESO9M14Yit7mzB1WDK2z9g3r+zLxENdJ5JG/ZskKe+TO4Diqi5OJt/h8yspS1ck8
LJtCole9919umByg5oruflqIlQIDAQABozUwMzALBgNVHQ8EBAMCAgQwDAYDVR0T
BAUwAwEB/zAWBgNVHREEDzANggtjYS5pc3Rpby5pbzANBgkqhkiG9w0BAQsFAAOC
AQEAltHEhhyAsve4K4bLgBXtHwWzo6SpFzdAfXpLShpOJNtQNERb3qg6iUGQdY+w
A2BpmSkKr3Rw/6ClP5+cCG7fGocPaZh+c+4Nxm9suMuZBZCtNOeYOMIfvCPcCS+8
PQ/0hC4/0J3WJKzGBssaaMufJxzgFPPtDJ998kY8rlROghdSaVt423/jXIAYnP3Y
05n8TGERBj7TLdtIVbtUIx3JHAo3PWJywA6mEDovFMJhJERp9sDHIr1BbhXK1TFN
Z6HNH6gInkSSMtvC4Ptejb749PTaePRPF7ID//eq/3AH8UK50F3TQcLjEqWUsJUn
aFKltOc+RAjzDklcUPeG4Y6eMA==
-----END CERTIFICATE-----`

var testRsaRootCert = `-----BEGIN CERTIFICATE-----
MIID7TCCAtWgAwIBAgIJAOIRDhOcxsx6MA0GCSqGSIb3DQEBCwUAMIGLMQswCQYD
VQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJU3Vubnl2YWxl
MQ4wDAYDVQQKDAVJc3RpbzENMAsGA1UECwwEVGVzdDEQMA4GA1UEAwwHUm9vdCBD
QTEiMCAGCSqGSIb3DQEJARYTdGVzdHJvb3RjYUBpc3Rpby5pbzAgFw0xODAxMjQx
OTE1NTFaGA8yMTE3MTIzMTE5MTU1MVowgYsxCzAJBgNVBAYTAlVTMRMwEQYDVQQI
DApDYWxpZm9ybmlhMRIwEAYDVQQHDAlTdW5ueXZhbGUxDjAMBgNVBAoMBUlzdGlv
MQ0wCwYDVQQLDARUZXN0MRAwDgYDVQQDDAdSb290IENBMSIwIAYJKoZIhvcNAQkB
FhN0ZXN0cm9vdGNhQGlzdGlvLmlvMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB
CgKCAQEA38uEfAatzQYqbaLou1nxJ348VyNzumYMmDDt5pbLYRrCo2pS3ki1ZVDN
8yxIENJFkpKw9UctTGdbNGuGCiSDP7uqF6BiVn+XKAU/3pnPFBbTd0S33NqbDEQu
IYraHSl/tSk5rARbC1DrQRdZ6nYD2KrapC4g0XbjY6Pu5l4y7KnFwSunnp9uqpZw
uERv/BgumJ5QlSeSeCmhnDhLxooG8w5tC2yVr1yDpsOHGimP/mc8Cds4V0zfIhQv
YzfIHphhE9DKjmnjBYLOdj4aycv44jHnOGc+wvA1Jqsl60t3wgms+zJTiWwABLdw
zgMAa7yxLyoV0+PiVQud6k+8ZoIFcwIDAQABo1AwTjAdBgNVHQ4EFgQUOUYGtUyh
euxO4lGe4Op1y8NVoagwHwYDVR0jBBgwFoAUOUYGtUyheuxO4lGe4Op1y8NVoagw
DAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEANXLyfAs7J9rmBamGJvPZ
ltx390WxzzLFQsBRAaH6rgeipBq3dR9qEjAwb6BTF+ROmtQzX+fjstCRrJxCto9W
tC8KvXTdRfIjfCCZjhtIOBKqRxE4KJV/RBfv9xD5lyjtCPCQl3Ia6MSf42N+abAK
WCdU6KCojA8WB9YhSCzza3aQbPTzd26OC/JblJpVgtus5f8ILzCsz+pbMimgTkhy
AuhYRppJaQ24APijsEC9+GIaVKPg5IwWroiPoj+QXNpshuvqVQQXvGaRiq4zoSnx
xAJz+w8tjrDWcf826VN14IL+/Cmqlg/rIfB5CHdwVIfWwpuGB66q/UiPegZMNs8a
3g==
-----END CERTIFICATE-----`

var testRsaCertChain = `-----BEGIN CERTIFICATE-----
MIIDnzCCAoegAwIBAgIJAON1ifrBZ2/BMA0GCSqGSIb3DQEBCwUAMIGLMQswCQYD
VQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJU3Vubnl2YWxl
MQ4wDAYDVQQKDAVJc3RpbzENMAsGA1UECwwEVGVzdDEQMA4GA1UEAwwHUm9vdCBD
QTEiMCAGCSqGSIb3DQEJARYTdGVzdHJvb3RjYUBpc3Rpby5pbzAgFw0xODAxMjQx
OTE1NTFaGA8yMTE3MTIzMTE5MTU1MVowWTELMAkGA1UEBhMCVVMxEzARBgNVBAgT
CkNhbGlmb3JuaWExEjAQBgNVBAcTCVN1bm55dmFsZTEOMAwGA1UEChMFSXN0aW8x
ETAPBgNVBAMTCElzdGlvIENBMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC
AQEAyzCxr/xu0zy5rVBiso9ffgl00bRKvB/HF4AX9/ytmZ6Hqsy13XIQk8/u/By9
iCvVwXIMvyT0CbiJq/aPEj5mJUy0lzbrUs13oneXqrPXf7ir3HzdRw+SBhXlsh9z
APZJXcF93DJU3GabPKwBvGJ0IVMJPIFCuDIPwW4kFAI7R/8A5LSdPrFx6EyMXl7K
M8jekC0y9DnTj83/fY72WcWX7YTpgZeBHAeeQOPTZ2KYbFal2gLsar69PgFS0Tom
ESO9M14Yit7mzB1WDK2z9g3r+zLxENdJ5JG/ZskKe+TO4Diqi5OJt/h8yspS1ck8
LJtCole9919umByg5oruflqIlQIDAQABozUwMzALBgNVHQ8EBAMCAgQwDAYDVR0T
BAUwAwEB/zAWBgNVHREEDzANggtjYS5pc3Rpby5pbzANBgkqhkiG9w0BAQsFAAOC
AQEAltHEhhyAsve4K4bLgBXtHwWzo6SpFzdAfXpLShpOJNtQNERb3qg6iUGQdY+w
A2BpmSkKr3Rw/6ClP5+cCG7fGocPaZh+c+4Nxm9suMuZBZCtNOeYOMIfvCPcCS+8
PQ/0hC4/0J3WJKzGBssaaMufJxzgFPPtDJ998kY8rlROghdSaVt423/jXIAYnP3Y
05n8TGERBj7TLdtIVbtUIx3JHAo3PWJywA6mEDovFMJhJERp9sDHIr1BbhXK1TFN
Z6HNH6gInkSSMtvC4Ptejb749PTaePRPF7ID//eq/3AH8UK50F3TQcLjEqWUsJUn
aFKltOc+RAjzDklcUPeG4Y6eMA==
-----END CERTIFICATE-----
`

func GetKubeConfig() *rest.Config {
	if kubeConfig != nil {
		return kubeConfig
	}
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	kubeConfig = cfg
	return cfg
}

func GetKubeClient() *kubernetes.Clientset {
	if kubeClient != nil {
		return kubeClient
	}
	cfg := GetKubeConfig()
	client, err := kubernetes.NewForConfig(cfg)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	kubeClient = client
	return client
}

func GetApiExtsClient() apiexts.Interface {
	if apiExtsClient != nil {
		return apiExtsClient
	}
	cfg := GetKubeConfig()
	client, err := apiexts.NewForConfig(cfg)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	apiExtsClient = client
	return client
}

func GetSecurityClient() *security.Clientset {
	securityClient, err := security.NewForConfig(GetKubeConfig())
	ExpectWithOffset(1, err).To(BeNil())
	return securityClient
}

func GetSecretClient() istiosecret.IstioCacertsSecretClient {
	if secretClient != nil {
		return secretClient
	}
	client, err := istiosecret.NewIstioCacertsSecretClient(&factory.KubeSecretClientFactory{
		Clientset: GetKubeClient(),
	})
	ExpectWithOffset(1, err).Should(BeNil())
	err = client.Register()
	ExpectWithOffset(1, err).Should(BeNil())
	secretClient = client
	return secretClient
}

func TryCreateNamespace(namespace string) {
	client := GetKubeClient()
	resource := &kubecore.Namespace{
		ObjectMeta: kubemeta.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := client.CoreV1().Namespaces().Create(resource)
	if err != nil {
		ExpectWithOffset(1, apierrors.IsAlreadyExists(err)).To(BeTrue())
	}
}

func TerminateNamespace(namespace string) {
	client := GetKubeClient()
	gracePeriod := int64(0)
	deleteOptions := &kubemeta.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
	}
	client.CoreV1().Pods(namespace).DeleteCollection(deleteOptions, kubemeta.ListOptions{})
	client.CoreV1().Namespaces().Delete(namespace, deleteOptions)
}

func TerminateNamespaceBlocking(namespace string) {
	TerminateNamespace(namespace)
	WaitForTerminatedNamespace(namespace)
}

func WaitForTerminatedNamespace(namespace string) {
	client := GetKubeClient()
	EventuallyWithOffset(1, func() error {
		_, err := client.CoreV1().Namespaces().Get(namespace, kubemeta.GetOptions{})
		return err
	}, "120s", "1s").ShouldNot(BeNil()) // will be non-nil when NS is gone
}

func WaitForAvailablePodsWithTimeout(namespace string, timeout string) int {
	// use helper function so that stack offset is consistent
	return waitForAvailablePodsWithTimeout(namespace, timeout)
}

func waitForAvailablePodsWithTimeout(namespace, timeout string) int {
	client := GetKubeClient()
	var podNum int

	EventuallyWithOffset(2, func() (bool, error) {
		podList, err := client.CoreV1().Pods(namespace).List(kubemeta.ListOptions{})
		if err != nil {
			return false, err
		}
		podNum = len(podList.Items)
		done := true
		for _, pod := range podList.Items {
			for _, condition := range pod.Status.Conditions {
				if pod.Status.Phase == kubecore.PodSucceeded {
					continue
				}
				if condition.Type == kubecore.PodReady && condition.Status != kubecore.ConditionTrue {
					done = false
				}
			}
		}
		return done, nil
	}, timeout, "1s").Should(BeTrue())
	return podNum
}

func WaitForDeletedPodsWithTimeout(namespace string, timeout string) {
	client := GetKubeClient()
	Eventually(func() bool {
		podList, err := client.CoreV1().Pods(namespace).List(kubemeta.ListOptions{})
		Expect(err).To(BeNil())
		return len(podList.Items) == 0
	}, timeout, "1s").Should(BeTrue())
}

func WaitForAvailablePods(namespace string) int {
	return waitForAvailablePodsWithTimeout(namespace, "120s")
}

func WaitForDeletedPods(namespace string) {
	WaitForDeletedPodsWithTimeout(namespace, "120s")
}

func GetMeshClient(kubeCache *kube.KubeCache) v1.MeshClient {
	meshClient, err := v1.NewMeshClient(&factory.KubeResourceClientFactory{
		Crd:         v1.MeshCrd,
		Cfg:         GetKubeConfig(),
		SharedCache: kubeCache,
	})
	ExpectWithOffset(1, err).Should(BeNil())
	err = meshClient.Register()
	ExpectWithOffset(1, err).Should(BeNil())
	return meshClient
}

func GetUpstreamClient(kubeCache *kube.KubeCache) gloo.UpstreamClient {
	if upstreamClient != nil {
		return upstreamClient
	}
	client, err := gloo.NewUpstreamClient(&factory.KubeResourceClientFactory{
		Crd:         gloo.UpstreamCrd,
		Cfg:         GetKubeConfig(),
		SharedCache: kubeCache,
	})
	ExpectWithOffset(1, err).Should(BeNil())
	err = client.Register()
	ExpectWithOffset(1, err).Should(BeNil())
	upstreamClient = client
	return upstreamClient
}

func DeleteCrb(crbName string) {
	client := GetKubeClient()
	client.RbacV1().ClusterRoleBindings().Delete(crbName, &kubemeta.DeleteOptions{})
}

func CrbDoesntExist(crbName string) bool {
	client := GetKubeClient()
	_, err := client.RbacV1().ClusterRoleBindings().Get(crbName, kubemeta.GetOptions{})
	return apierrors.IsNotFound(err)
}

func DeleteWebhookConfigIfExists(webhookName string) {
	client := GetKubeClient()
	hooks, err := client.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().List(kubemeta.ListOptions{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	for _, hook := range hooks.Items {
		if strings.HasSuffix(hook.Name, webhookName) {
			client.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Delete(hook.Name, &kubemeta.DeleteOptions{})
		}
	}
}

func WebhookConfigNotFound(webhookName string) bool {
	client := GetKubeClient()
	_, err := client.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Get(webhookName, kubemeta.GetOptions{})
	return apierrors.IsNotFound(err)
}

func GetConsulServerPodName(namespace string) string {
	client := GetKubeClient()
	podList, err := client.CoreV1().Pods(namespace).List(kubemeta.ListOptions{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	for _, pod := range podList.Items {
		if strings.Contains(pod.Name, "consul-mesh-server-0") {
			return pod.Name
		}
	}
	// Should not have happened
	ExpectWithOffset(1, false).To(BeTrue())
	return ""
}

// New creates a new and initialized tunnel.
func CreateConsulTunnel(namespace string, port int) (*helmkube.Tunnel, error) {
	podName := GetConsulServerPodName(namespace)
	t := helmkube.NewTunnel(GetKubeClient().CoreV1().RESTClient(), GetKubeConfig(), namespace, podName, port)
	return t, t.ForwardPort()
}

func CreateTestEcSecret(namespace string, name string) (*istiosecret.IstioCacertsSecret, *core.ResourceRef) {
	secret := &istiosecret.IstioCacertsSecret{
		Metadata: core.Metadata{
			Namespace: namespace,
			Name:      name,
		},
		CaCert:    TestEcRoot,
		CaKey:     testEcKey,
		RootCert:  TestEcRoot,
		CertChain: testCertChain,
	}
	GetSecretClient().Delete(namespace, name, clients.DeleteOpts{})
	_, err := GetSecretClient().Write(secret, clients.WriteOpts{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ref := &core.ResourceRef{
		Namespace: namespace,
		Name:      name,
	}
	return secret, ref
}

func CreateTestRsaSecret(namespace string, name string) (*istiosecret.IstioCacertsSecret, *core.ResourceRef) {
	secret := &istiosecret.IstioCacertsSecret{
		Metadata: core.Metadata{
			Namespace: namespace,
			Name:      name,
		},
		CaCert:    testRsaCaCert,
		CaKey:     testRsaCaKey,
		RootCert:  testRsaRootCert,
		CertChain: testCertChain,
	}
	return createSecret(namespace, name, secret)
}

func createSecret(namespace string, name string, secret *istiosecret.IstioCacertsSecret) (*istiosecret.IstioCacertsSecret, *core.ResourceRef) {
	GetSecretClient().Delete(namespace, name, clients.DeleteOpts{})
	_, err := GetSecretClient().Write(secret, clients.WriteOpts{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ref := &core.ResourceRef{
		Namespace: namespace,
		Name:      name,
	}
	return secret, ref
}

func CheckCertMatchesConsul(consulTunnelPort int, rootCert string) {
	config := &api.Config{
		Address: fmt.Sprintf("127.0.0.1:%d", consulTunnelPort),
	}
	client, err := api.NewClient(config)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	var queryOpts api.QueryOptions
	currentConfig, _, err := client.Connect().CAGetConfig(&queryOpts)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	currentRoot := currentConfig.Config["RootCert"]
	ExpectWithOffset(1, currentRoot).To(BeEquivalentTo(rootCert))
}

func CheckCertMatchesIstio(installNamespace string) {
	actual, err := GetSecretClient().Read(installNamespace, secret.CustomRootCertificateSecretName, clients.ReadOpts{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ExpectWithOffset(1, actual.RootCert).Should(BeEquivalentTo(testRsaRootCert))
	ExpectWithOffset(1, actual.CaCert).Should(BeEquivalentTo(testRsaCaCert))
	ExpectWithOffset(1, actual.CaKey).Should(BeEquivalentTo(testRsaCaKey))
	ExpectWithOffset(1, actual.CertChain).Should(BeEquivalentTo(testCertChain))
}

func UninstallHelmRelease(releaseName string) error {
	// helm install
	helmClient, err := helm.GetHelmClient(context.TODO())
	if err != nil {
		return err
	}
	_, err = helmClient.DeleteRelease(releaseName, helmlib.DeletePurge(true))
	helm.Teardown()
	return err
}

func HelmReleaseDoesntExist(releaseName string) bool {
	helmClient, err := helm.GetHelmClient(context.TODO())
	if err != nil {
		return false
	}
	statuses := []release.Status_Code{
		release.Status_UNKNOWN,
		release.Status_DEPLOYED,
		release.Status_DELETED,
		release.Status_SUPERSEDED,
		release.Status_FAILED,
		release.Status_DELETING,
		release.Status_PENDING_INSTALL,
		release.Status_PENDING_UPGRADE,
		release.Status_PENDING_ROLLBACK,
	}
	// equivalent to "--all" option
	list, err := helmClient.ListReleases(helmlib.ReleaseListStatuses(statuses))
	if err != nil {
		return false
	}
	// No releases == successfully deleted
	if list == nil {
		return true
	}
	for _, item := range list.Releases {
		if item.Name == releaseName {
			return false
		}
	}
	return true
}

func TryDeleteIstioCrds() {
	crdClient := GetApiExtsClient()
	crdList, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().List(kubemeta.ListOptions{})
	if err != nil {
		return
	}
	for _, crd := range crdList.Items {
		//TODO: use labels
		if strings.Contains(crd.Name, "istio.io") {
			crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, &kubemeta.DeleteOptions{})
		}
	}
}

func IstioCrdsDontExist() bool {
	crdClient := GetApiExtsClient()
	crdList, err := crdClient.ApiextensionsV1beta1().CustomResourceDefinitions().List(kubemeta.ListOptions{})
	if err != nil {
		return false
	}
	for _, crd := range crdList.Items {
		if strings.Contains(crd.Name, "istio.io") {
			return false
		}
	}
	return true
}

func GetUpstreamNames(upstreamClient gloo.UpstreamClient) ([]string, error) {
	ul, err := upstreamClient.List("gloo-system", clients.ListOpts{})
	if err != nil {
		return nil, err
	}
	ups := []string{}
	for _, up := range ul {
		ups = append(ups, up.Metadata.Name)
	}
	return ups, nil
}
