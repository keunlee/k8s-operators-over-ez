package task001

import (
	task001v1alpha1 "github.com/keunlee/task-001-operator/pkg/apis/task001/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestAdd(t *testing.T) {
	type args struct {
		mgr manager.Manager
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Add(tt.args.mgr); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReconcileTask001_Reconcile(t *testing.T) {
	type fields struct {
		client client.Client
		scheme *runtime.Scheme
	}
	type args struct {
		request reconcile.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    reconcile.Result
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReconcileTask001{
				client: tt.fields.client,
				scheme: tt.fields.scheme,
			}
			got, err := r.Reconcile(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reconcile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_add(t *testing.T) {
	type args struct {
		mgr manager.Manager
		r   reconcile.Reconciler
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := add(tt.args.mgr, tt.args.r); (err != nil) != tt.wantErr {
				t.Errorf("add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newPodForCR(t *testing.T) {
	type args struct {
		cr *task001v1alpha1.Task001
	}
	tests := []struct {
		name string
		args args
		want *corev1.Pod
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newPodForCR(tt.args.cr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newPodForCR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newReconciler(t *testing.T) {
	type args struct {
		mgr manager.Manager
	}
	tests := []struct {
		name string
		args args
		want reconcile.Reconciler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newReconciler(tt.args.mgr); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newReconciler() = %v, want %v", got, tt.want)
			}
		})
	}
}