package controllers

import (
	"time"

	"k8s.io/utils/pointer"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/types"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appsv1 "k8s.io/api/apps/v1"

	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testRegistryManager1 = &TestRegistryManager{}

var _ = Describe("Deployment Controller Test", func() {

	const (
		DeploymentName                = "test-deployment"
		DeploymentNamespace           = "ns1"
		DestinationRegistrySecretName = "destination-registry-creds"
		timeout                       = time.Second * 10
		interval                      = time.Millisecond * 250
	)

	var actualSrcImageNames, actualDstImageNames []string
	var actualSrcRegistryCredentials []*RegistryCredentials
	var actualDstRegistryCredentials *RegistryCredentials

	testRegistryManager1.copyImageStub = func(srcImage, dstImage string, srcRegistryCredentials, dstRegistryCredentials *RegistryCredentials) {
		actualSrcImageNames = append(actualSrcImageNames, srcImage)
		actualDstImageNames = append(actualDstImageNames, dstImage)
		actualSrcRegistryCredentials = append(actualSrcRegistryCredentials, srcRegistryCredentials)
		actualDstRegistryCredentials = dstRegistryCredentials

	}

	Context("Update image with Backup registry", func() {
		It("Should create successfully", func() {

			ns := &corev1.Namespace{}
			ns.Name = DeploymentNamespace
			Expect(k8sClient.Create(context.Background(), ns)).Should(Succeed())

			secret1 := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret1-deployment",
					Namespace: DeploymentNamespace,
				},
				Type: "kubernetes.io/dockerconfigjson",
				Data: map[string][]byte{
					".dockerconfigjson": SrcRegAuth1,
				},
			}
			Expect(k8sClient.Create(context.Background(), secret1)).Should(Succeed())

			secret2 := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "secret2-deployment",
					Namespace: DeploymentNamespace,
				},
				Type: "kubernetes.io/dockerconfigjson",
				Data: map[string][]byte{
					".dockerconfigjson": SrcRegAuth2,
				},
			}
			Expect(k8sClient.Create(context.Background(), secret2)).Should(Succeed())

			// Create
			deployment := &appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      DeploymentName,
					Namespace: DeploymentNamespace,
				},
				Spec: appsv1.DeploymentSpec{
					Replicas: pointer.Int32Ptr(1),
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{
							"app": "test-cronjob",
						},
					},
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{
							Labels: map[string]string{
								"app": "test-cronjob",
							},
						},
						Spec: corev1.PodSpec{
							ImagePullSecrets: []corev1.LocalObjectReference{{Name: "secret1-deployment"}, {Name: "secret2-deployment"}},
							Containers: []corev1.Container{
								{
									Name:  "test-cont1",
									Image: SrcImageNames[0],
								},
								{
									Name:  "test-cont2",
									Image: SrcImageNames[1],
								},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(context.Background(), deployment)).Should(Succeed())

			deploymentLookupKey := types.NamespacedName{Name: DeploymentName, Namespace: DeploymentNamespace}
			createdDeployment := &appsv1.Deployment{}

			By("Expecting image to be updated in containers")
			Eventually(func() ([]string, error) {
				err := k8sClient.Get(context.Background(), deploymentLookupKey, createdDeployment)
				if err != nil {
					return nil, err
				}

				var names []string
				for _, container := range createdDeployment.Spec.Template.Spec.Containers {
					names = append(names, container.Image)
				}
				return names, nil
			}, timeout, interval).Should(ConsistOf(DstImageNames), "should list updated image name in container list", SrcImageNames)

			By("Expecting destination secret to be created")
			Eventually(func() (string, error) {
				secret := &corev1.Secret{}
				err := k8sClient.Get(context.Background(), types.NamespacedName{Name: DestinationRegistrySecretName, Namespace: DeploymentNamespace}, secret)
				if err != nil {
					return "", err
				}

				return string(secret.Data[".dockerconfigjson"]), nil
			}, timeout, interval).Should(MatchJSON(DstRegAuth), "should list destination secret data")

			By("Expecting src image name in copy image")
			Eventually(func() ([]string, error) {
				return actualSrcImageNames, nil
			}, timeout, interval).Should(Equal(SrcImageNames), "should list src image name in copy method")

			By("Expecting dst image name in copy image")
			Eventually(func() ([]string, error) {
				return actualDstImageNames, nil
			}, timeout, interval).Should(Equal(DstImageNames), "should list dst image name in copy method")

			By("Expecting dst image credential in copy image")
			Eventually(func() (*RegistryCredentials, error) {
				return actualDstRegistryCredentials, nil
			}, timeout, interval).Should(Equal(DstRegistryCredentials), "should list dst registry credentials in copy method")

			By("Expecting src image credential in copy image")
			Eventually(func() ([]*RegistryCredentials, error) {
				return actualSrcRegistryCredentials, nil
			}, timeout, interval).Should(Equal(SrcRegistryCredentialList), "should list src registry credentials in copy method")

		})
	})
})
