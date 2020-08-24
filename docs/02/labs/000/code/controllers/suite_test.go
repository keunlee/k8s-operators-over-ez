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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"

	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	operatorsoverezv1alpha1 "github.com/mydomain/operators-over-ez/api/v1alpha1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

var k8sManager ctrl.Manager
var opsOverEasyReconciler *OpsOverEasyReconciler
var crdInstance = &operatorsoverezv1alpha1.OpsOverEasy{}
var testCtx = context.Background()
var crKey = types.NamespacedName{
	Name:      "operator-overeasy",
	Namespace: "default",
}

func getCrd(withSpecification bool) *operatorsoverezv1alpha1.OpsOverEasy {
	var crd *operatorsoverezv1alpha1.OpsOverEasy

	if withSpecification {
		crd = &operatorsoverezv1alpha1.OpsOverEasy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crKey.Name,
				Namespace: crKey.Namespace,
			},

			Spec: operatorsoverezv1alpha1.OpsOverEasySpec{
				Timeout: 30,
				Message: "message",
			},
		}
	} else {
		crd = &operatorsoverezv1alpha1.OpsOverEasy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      crKey.Name,
				Namespace: crKey.Namespace,
			},
		}
	}

	return crd
}

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = operatorsoverezv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	k8sManager, err = ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sManager).ToNot(BeNil())

	opsOverEasyReconciler = &OpsOverEasyReconciler{
		Client: k8sClient,
		Log:    ctrl.Log.WithName("controllers").WithName("OpsOverEasy"),
		Scheme: scheme.Scheme,
	}
	err = (opsOverEasyReconciler).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	close(done)
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})
