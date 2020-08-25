/*


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
	"fmt"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	operatorsoverezv1alpha1 "github.com/mydomain/operators-over-ez/api/v1alpha1"
)

// OpsOverEasyReconciler reconciles a OpsOverEasy object
type OpsOverEasyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=operators-over-ez.mydomain.com,resources=opsovereasies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operators-over-ez.mydomain.com,resources=opsovereasies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *OpsOverEasyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("opsovereasy", req.NamespacedName)

	// your logic here
	instance := &operatorsoverezv1alpha1.OpsOverEasy{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	//resp, err := http.Get("http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response")
	//if ( resp != nil) {}

	pod := newPodForCR(instance)

	// Set Mycrd instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err == nil {
		if found.Status.Phase == corev1.PodSucceeded {
			instance.Status.MessageLogged = true
			instance.Status.TimeoutExpired = true

			err = r.Client.Status().Update(context.TODO(), instance)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.Client.Create(context.TODO(), pod)
		if err != nil {
			return ctrl.Result{}, err
		}

		// Pod created successfully - don't requeue
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// Pod already exists - don't requeue
	r.Log.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return ctrl.Result{}, nil
}

func (r *OpsOverEasyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorsoverezv1alpha1.OpsOverEasy{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *operatorsoverezv1alpha1.OpsOverEasy) *corev1.Pod {
	timeout := cr.Spec.Timeout
	message := cr.Spec.Message

	labels := map[string]string{
		"app": cr.Name,
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "busybox",
					Image: "busybox",
					Args:  []string{"/bin/sh", "-c", fmt.Sprintf("sleep %d; echo '%s'", timeout, message)},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}
