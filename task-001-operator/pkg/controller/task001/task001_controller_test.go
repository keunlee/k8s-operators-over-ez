package task001

import (
	"testing"

	task001v1alpha1 "github.com/keunlee/task-001-operator/pkg/apis/task001/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestReconcileTask001_Reconcile(t *testing.T) {
	type fields struct {
		client client.Client
		scheme *runtime.Scheme
	}

	type args struct {
		request reconcile.Request
	}

	var (
		name            = "task001-operator"
		namespace       = "task001"
	)

	// A Task001 object with metadata and spec.
	task001 := &task001v1alpha1.Task001{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: task001v1alpha1.Task001Spec{
		},
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{ task001 }

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(task001v1alpha1.SchemeGroupVersion, task001)

	// Create a fake client to mock API calls.
	cl := fake.NewFakeClientWithScheme(s, objs...)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    reconcile.Result
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "",
			fields:  fields{
				client: cl,
				scheme: s,
			},
			args:    args{},
			want:    reconcile.Result{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create a ReconcileTask001 object with the scheme and fake client.
			r := &ReconcileTask001{
				client: tt.fields.client,
				scheme: tt.fields.scheme,
			}

			// Mock request to simulate Reconcile() being called on an event for a
			// watched resource .
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      name,
					Namespace: namespace,
				},
			}

			got, err := r.Reconcile(req)

			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
				return
			}

			// Check the result of reconciliation to make sure it has the desired state.
			if got.Requeue {
				t.Error("reconcile unexpected requeued request")
				return
			}
		})
	}
}