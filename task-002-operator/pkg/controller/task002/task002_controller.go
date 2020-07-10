package task002

import (
	"context"
	"reflect"

	task002v1alpha1 "github.com/keunlee/task-002-operator/pkg/apis/task002/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_task002")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Task002 Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileTask002{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("task002-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Task002
	err = c.Watch(&source.Kind{Type: &task002v1alpha1.Task002{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Task002
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &task002v1alpha1.Task002{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileTask002 implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileTask002{}

// ReconcileTask002 reconciles a Task002 object
type ReconcileTask002 struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Task002 object and makes changes based on the state read
// and what is in the Task002.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileTask002) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Task002")

	// Fetch the Task002 instance
	instance := &task002v1alpha1.Task002{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Retrieve a list of existing pods in our namespace
	// See: https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/client#example-Client--List
	task002 := instance
	podList := &corev1.PodList{}

	// // basic pod list retrieval
	// if err = r.client.List(context.Background(), podList); err != nil {
	// 	return reconcile.Result{}, err
	// }

	// pod list retreival, using selectors to filter pods
	// See: https://godoc.org/k8s.io/apimachinery/pkg/labels#SelectorFromSet
	lbs := map[string]string{
		"app":     task002.Name,
		"version": "v0.1",
		"crType":  "Task002",
	}
	labelSelector := labels.SelectorFromSet(lbs)
	listOps := &client.ListOptions{Namespace: task002.Namespace, LabelSelector: labelSelector}
	if err = r.client.List(context.TODO(), podList, listOps); err != nil {
		return reconcile.Result{}, err
	}

	// Count the pods that are pending or running as available
	var available []corev1.Pod
	for _, pod := range podList.Items {
		if pod.ObjectMeta.DeletionTimestamp != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodPending {
			available = append(available, pod)
		}
	}
	numAvailable := int(len(available))
	availableNames := []string{}
	for _, pod := range available {
		availableNames = append(availableNames, pod.ObjectMeta.Name)
	}

	// Update the status if necessary
	status := task002v1alpha1.Task002Status{
		ListedPods: availableNames,
	}
	if !reflect.DeepEqual(task002.Status, status) {
		task002.Status = status

		// Update the status
		// See: https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/client#StatusWriter
		err = r.client.Status().Update(context.TODO(), task002)
		if err != nil {
			reqLogger.Error(err, "Failed to update Task002 status")
			return reconcile.Result{}, err
		}
	}

	// scale pods down to 'NumberOfPods' specification
	if numAvailable > task002.Spec.NumberOfPods {
		reqLogger.Info("Scaling down pods", "Currently available", numAvailable, "Required replicas", task002.Spec.NumberOfPods)
		diff := numAvailable - task002.Spec.NumberOfPods
		dpods := available[:diff]
		for _, dpod := range dpods {

			// Delete specified pod
			// See: https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/client#example-Client--Delete
			err = r.client.Delete(context.TODO(), &dpod)
			if err != nil {
				reqLogger.Error(err, "Failed to delete pod", "pod.name", dpod.Name)
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{Requeue: true}, nil
	}

	// scale pods up to 'NumberOfPods' specification
	if numAvailable < task002.Spec.NumberOfPods {
		reqLogger.Info("Scaling up pods", "Currently available", numAvailable, "Required replicas", task002.Spec.NumberOfPods)
		// Define a new Pod object
		pod := newPodForCR(task002)

		// Set Task002 instance as the owner and controller
		// See: https://godoc.org/sigs.k8s.io/controller-runtime/pkg/controller/controllerutil#SetControllerReference
		if err := controllerutil.SetControllerReference(task002, pod, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Create specified pod
		// See: https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg/client#example-Client--Create
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			reqLogger.Error(err, "Failed to create pod", "pod.name", pod.Name)
			return reconcile.Result{}, err
		}
		return reconcile.Result{Requeue: true}, nil
	}

	reqLogger.Info("msg", "bp")

	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *task002v1alpha1.Task002) *corev1.Pod {
	labels := map[string]string{
		"app":     cr.Name,
		"version": "v0.1",
		"crType":  "Task002",
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: cr.Name + "-pod", // use 'GenerateName' instead of 'Name', to create a unique name when this resource is created
			Namespace:    cr.Namespace,
			Labels:       labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}
