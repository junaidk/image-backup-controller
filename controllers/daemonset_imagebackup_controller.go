/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DaemonsetImageBackupReconciler reconciles a ImageBackup object
type DaemonsetImageBackupReconciler struct {
	client.Client
	Scheme                    *runtime.Scheme
	RegistryManager           RegistryManager
	BackUpRegistryCredentials *RegistryCredentials
	IgnoreNamespaces          []string
}

//+kubebuilder:rbac:groups=apps,resources=daemonset,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=daemonset/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=apps,resources=daemonset/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the ImageBackup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *DaemonsetImageBackupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lg := log.FromContext(ctx)

	daemonset := &appsv1.DaemonSet{}

	err := r.Client.Get(ctx, req.NamespacedName, daemonset)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	for _, container := range daemonset.Spec.Template.Spec.Containers {
		lg.Info("Image", "namespace", daemonset.Namespace, "name", daemonset.Name, "image", container.Image)
	}
	var srcImages []string
	var dstImages []string

	// get src and dst image name list
	for _, container := range daemonset.Spec.Template.Spec.Containers {
		srcImages = append(srcImages, container.Image)
		dstImages = append(dstImages, getDestinationImageName(container.Image, r.BackUpRegistryCredentials.URL, r.BackUpRegistryCredentials.Username))
	}

	// get registry credentials from ImagePullSecrets
	srcRegistryCredentials, err := getRegistryCredentials(ctx, r.Client, daemonset.Spec.Template.Spec.ImagePullSecrets, daemonset.Namespace)
	if err != nil {
		lg.Error(err, "failed to get registry credentials")
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	lg.Info("src", "creds", srcRegistryCredentials)

	// create destination registry secret
	dstRegistryDockerSecret, err := getDockerConfigSecret(r.BackUpRegistryCredentials.Username, r.BackUpRegistryCredentials.Password, r.BackUpRegistryCredentials.URL)
	if err != nil {
		lg.Error(err, "failed to get docker config secret")
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}
	if dstRegistryDockerSecret != nil {
		dstRegistryDockerSecret.Namespace = daemonset.Namespace
		err = createRegistrySecret(ctx, r.Client, dstRegistryDockerSecret)
		if err != nil {
			lg.Error(err, "failed to create registry secret")
			return ctrl.Result{RequeueAfter: time.Second * 10}, nil
		}
	}

	// copy images from src to dst.
	// TODO improvement. make image copy concurrent for multiple images (using go routines)
	for i, container := range daemonset.Spec.Template.Spec.Containers {
		if strings.Contains(container.Image, r.BackUpRegistryCredentials.URL) {
			continue
		}
		srcRegistryURL := strings.Split(container.Image, "/")[0]
		if true {
			srcRegistryCredential := &RegistryCredentials{}

			if len(srcRegistryCredentials) != 0 {
				if _, ok := srcRegistryCredentials[srcRegistryURL]; !ok {
					srcRegistryURL = DEFAULT_DOCKER_REGISTRY
				}
				srcRegistryCredential = srcRegistryCredentials[srcRegistryURL]
			}
			err := r.RegistryManager.CopyImage(ctx, srcImages[i], dstImages[i], srcRegistryCredential, r.BackUpRegistryCredentials)
			if err != nil {
				lg.Error(err, "failed to copy image")
				return ctrl.Result{RequeueAfter: time.Second * 10}, nil
			}
		}
	}

	// update image name in daemonset
	for i := range daemonset.Spec.Template.Spec.Containers {
		if true {
			daemonset.Spec.Template.Spec.Containers[i].Image = dstImages[i]
		}
	}

	if dstRegistryDockerSecret != nil {
		daemonset.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{{Name: dstRegistryDockerSecret.Name}}
	}

	err = r.Client.Update(ctx, daemonset)
	if err != nil {
		lg.Error(err, "failed to update daemonset")
		return ctrl.Result{RequeueAfter: time.Second * 10}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DaemonsetImageBackupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
		WithEventFilter(ignorePredicate(r.IgnoreNamespaces)).
		Complete(r)
}
