package controllers

import (
	"context"
	"fmt"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	DEFAULT_DOCKER_REGISTRY = "index.docker.io"
)

func getDestinationImageName(image, registryURL, registryUser string) string {
	imageName := strings.Split(image, "/")
	dstImage := registryURL + "/" + registryUser + "/" + imageName[len(imageName)-1]

	return dstImage
}

func getRegistrySecret(ctx context.Context, k8sClient client.Client, name, namespace string) (*corev1.Secret, error) {
	regSecret := &corev1.Secret{}
	regSecretName := name
	err := k8sClient.Get(ctx, client.ObjectKey{Name: regSecretName, Namespace: namespace}, regSecret)
	if err != nil {
		return nil, err
	}
	return regSecret, nil
}

func createRegistrySecret(ctx context.Context, k8sClient client.Client, secret *corev1.Secret) error {
	existingSecret := &corev1.Secret{}
	err := k8sClient.Get(ctx, client.ObjectKey{Name: secret.Name, Namespace: secret.Namespace}, existingSecret)
	if err != nil {
		if errors.IsNotFound(err) {
			err = k8sClient.Create(ctx, secret)
			if err != nil {
				return err
			}
		} else {
			existingSecret.Data = secret.Data
			err = k8sClient.Update(ctx, existingSecret)
			return err
		}
	}
	return nil
}

func getRegistryCredentials(ctx context.Context, k8sclient client.Client, imgPullSecrets []corev1.LocalObjectReference, namespace string) (map[string]*RegistryCredentials, error) {
	registryCredentials := make(map[string]*RegistryCredentials)

	for _, secret := range imgPullSecrets {
		registryCredential, err := getRegistryCredential(ctx, k8sclient, secret.Name, namespace)
		if err != nil {
			return nil, err
		}
		registryCredentials[registryCredential.URL] = &RegistryCredentials{
			URL:      registryCredential.URL,
			Username: registryCredential.Username,
			Password: registryCredential.Password,
		}
	}
	return registryCredentials, nil
}

func getRegistryCredential(ctx context.Context, k8sclient client.Client, secretName, namespace string) (*RegistryCredentials, error) {
	var regCreds *corev1.Secret
	var err error

	regCreds, err = getRegistrySecret(ctx, k8sclient, secretName, namespace)
	if err != nil {
		return nil, fmt.Errorf("error getting registry secret: %v", err)
	}

	var registryCreds *RegistryCredentials
	if regCreds != nil {
		if regCreds.Type == corev1.SecretTypeDockerConfigJson {
			registryCreds, err = getRegistryCredentialsFromDockerCfg(regCreds.Data[corev1.DockerConfigJsonKey])
			if err != nil {
				return nil, fmt.Errorf("failed to get auth from docker config: %v", err)
			}
		}
	}
	return registryCreds, nil
}

func ignorePredicate(ignoreNamespaces []string) predicate.Predicate {

	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			if strings.Contains(strings.Join(ignoreNamespaces, ","), e.Object.GetNamespace()) {
				return false
			}
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			if strings.Contains(strings.Join(ignoreNamespaces, ","), e.ObjectNew.GetNamespace()) {
				return false
			}
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}
