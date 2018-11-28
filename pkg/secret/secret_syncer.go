package secret

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"k8s.io/client-go/kubernetes"

	istiov1 "github.com/solo-io/supergloo/pkg/api/external/istio/encryption/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/supergloo/pkg/api/v1"
	"github.com/solo-io/supergloo/pkg/translator/kube"
)

type SecretSyncer struct {
	IstioNamespace string
	Kube           kubernetes.Interface
	SecretClient   istiov1.IstioCacertsSecretClient
	Preinstall     bool
}

const (
	CustomRootCertificateSecretName  = "cacerts"
	DefaultRootCertificateSecretName = "istio.default"
	istioLabelKey                    = "istio"
	citadelLabelValue                = "citadel"
)

func (s *SecretSyncer) SyncSecret(ctx context.Context, encryption *v1.Encryption, secretList istiov1.IstioCacertsSecretList) error {
	if encryption == nil {
		return nil
	}
	if !encryption.TlsEnabled {
		return nil
	}
	encryptionSecret := encryption.Secret
	if encryptionSecret == nil {
		return nil
	}
	sourceSecret, err := secretList.Find(encryptionSecret.Namespace, encryptionSecret.Name)
	if err != nil {
		return errors.Wrapf(err, "Error finding secret referenced in mesh config (%s:%s)",
			encryptionSecret.Namespace, encryptionSecret.Name)
	}
	// this is where custom root certs will live once configured, if not found existingSecret will be nil
	existingSecret, _ := secretList.Find(s.IstioNamespace, CustomRootCertificateSecretName)
	return s.syncSecret(ctx, sourceSecret, existingSecret)
}

func (s *SecretSyncer) syncSecret(ctx context.Context, sourceSecret, existingSecret *istiov1.IstioCacertsSecret) error {
	if err := validateTlsSecret(sourceSecret); err != nil {
		return errors.Wrapf(err, "invalid secret %v", sourceSecret.Metadata.Ref())
	}
	istioSecret := resources.Clone(sourceSecret).(*istiov1.IstioCacertsSecret)
	if existingSecret == nil {
		istioSecret.Metadata = core.Metadata{
			Namespace: s.IstioNamespace,
			Name:      CustomRootCertificateSecretName,
		}
		if _, err := s.SecretClient.Write(istioSecret, clients.WriteOpts{
			Ctx: ctx,
		}); err != nil {
			return errors.Wrapf(err, "creating tool tls secret %v for istio", istioSecret.Metadata.Ref())
		}
		return nil
	}

	// move secret over to destination name/namespace
	istioSecret.SetMetadata(existingSecret.Metadata)
	istioSecret.Metadata.Annotations["created_by"] = "supergloo"
	// nothing to do
	if istioSecret.Equal(existingSecret) {
		return nil
	}
	if _, err := s.SecretClient.Write(istioSecret, clients.WriteOpts{
		Ctx: ctx,
	}); err != nil {
		return errors.Wrapf(err, "updating tool tls secret %v for istio", istioSecret.Metadata.Ref())
	}
	return nil
}

func validateTlsSecret(secret *istiov1.IstioCacertsSecret) error {
	if secret.RootCert == "" {
		return errors.Errorf("Root cert is missing.")
	}
	if secret.CaKey == "" {
		return errors.Errorf("Private key is missing.")
	}
	return nil
}

func (s *SecretSyncer) deleteIstioDefaultSecret() error {
	// Using Kube API directly cause we don't expect this secret to be tagged and it should be mostly a one-time op
	return s.Kube.CoreV1().Secrets(s.IstioNamespace).Delete(DefaultRootCertificateSecretName, &metav1.DeleteOptions{})
}

func (s *SecretSyncer) restartCitadel() error {
	selector := make(map[string]string)
	selector[istioLabelKey] = citadelLabelValue
	return kube.RestartPods(s.Kube, s.IstioNamespace, selector)
}
