package controllers

import (
	"context"
	"github.com/go-logr/logr/testing"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	operatorsoverezv1alpha1 "github.com/mydomain/operators-over-ez/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var testCtx = context.Background()

var crKey = types.NamespacedName{
	Name:      "operator-overeasy",
	Namespace: "default",
}

var crdInstance = &operatorsoverezv1alpha1.OpsOverEasy{}

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

func getMockReconciler() *OpsOverEasyReconciler {
	// Register operator types with the runtime scheme.
	scheme := scheme.Scheme
	scheme.AddKnownTypes(operatorsoverezv1alpha1.GroupVersion, crdInstance)

	// Create a mock MycrdReconciler object with the scheme and client.
	reconciler := &OpsOverEasyReconciler{k8sClient, &testing.TestLogger{}, scheme}

	return reconciler
}

var _ = Describe("CR Controller", func() {

	const timeout = time.Second * 30
	const interval = time.Second * 1

	Context("BDD Test Scenarios", func() {
		Context("CR Instance with Specifications Provided", func() {
			BeforeEach(func() {
				// Given: An Operator Instance
				crdInstance = getCrd(true)
				err := k8sClient.Create(testCtx, crdInstance)
				Expect(err).ShouldNot(HaveOccurred())
			})

			AfterEach(func() {
				err := k8sClient.Delete(testCtx, crdInstance)
				Expect(err).ShouldNot(HaveOccurred())
			})

			//SCENARIO: Shutdown the busybox pod after a user specified amount of time in seconds
			//GIVEN: An Operator instance
			//WHEN: the specification timeout is set to a numeric value in seconds
			//THEN: the busy box pod will remain available for the specified timeout in seconds,
			//AND: shutdown after the specified amount timeout duration
			When("The specification timeout is set to a numeric value in seconds", func() {
				It("Should remain available for the specified timeout duration in seconds", func() {
					Expect(true).To(Equal(false))
				})

				It("Should shutdown after the specified amount timeout duration", func() {
					Expect(true).To(Equal(false))
				})
			})

			//SCENARIO: Log a user specified message before shutting down the busybox pod
			//GIVEN: An Operator instance
			//WHEN: the specification message is set to a string value
			//THEN: the busy box pod will log the message, from the message specification after the timeout duration has expired.
			When("The specification message is set to a string value", func() {
				It("Should log the message, from the message specification after the time out duration has expired", func() {
					Expect(true).To(Equal(false))
				})
			})

			//SCENARIO: Update status expired and logged when the busybox pod has expired
			//GIVEN: An Operator instance
			//WHEN: the busy box pod's duration has expired
			//THEN: set the expired status to true
			//AND: set the logged status to true
			When("The duration has expired", func() {
				It("Should set the expired status to true", func() {
					Expect(true).To(Equal(false))
				})
				It("Should set the logged status to true", func() {
					Expect(true).To(Equal(false))
				})
			})
		})

		Context("CR Instance with no Specifications Provided", func() {
			BeforeEach(func() {
				// Given: An Operator Instance
				crdInstance = getCrd(false)
				err := k8sClient.Create(testCtx, crdInstance)
				Expect(err).ShouldNot(HaveOccurred())
			})

			AfterEach(func() {
				err := k8sClient.Delete(testCtx, crdInstance)
				Expect(err).ShouldNot(HaveOccurred())
			})

			//SCENARIO: Retrieve the timeout and message from a given REST API if one and/or the other is not supplied.
			//GIVEN: An Operator instance
			//WHEN: the specification message OR timeout is NOT set
			//THEN: the busy box pod will supply these values from the following REST API: GET http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response
			When("The specification message OR timeout is NOT set", func() {
				It("Should supply these values from the following REST API: GET http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response", func() {
					Expect(true).To(Equal(false))
				})
			})
		})

	})

	Context("Unit Tests", func() {
		BeforeEach(func() {
			crdInstance = getCrd(true)
			err := k8sClient.Create(testCtx, crdInstance)
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			err := k8sClient.Delete(testCtx, crdInstance)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Should validate the CR was created", func() {
			By("Retrieving the CR successfully")
			fetched := &operatorsoverezv1alpha1.OpsOverEasy{}
			Expect(k8sClient.Get(testCtx, crKey, fetched)).Should(Succeed())

			By("Validating the expected CR specifications")
			Expect(fetched.Spec.Message).To(Equal("message"))
			Expect(fetched.Spec.Timeout).To(Equal(int32(30)))
		})

		It("Should reconcile the CR successfully", func() {
			podKey := types.NamespacedName{
				Namespace: "default",
				Name:      "operator-overeasy-pod",
			}

			// Create a mock reconciler object with the scheme and client.
			By("Leveraging an instance of the Reconciler")
			reconciler := getMockReconciler()

			// Mock request to simulate Reconcile() being called on an event for a watched resource .
			By("Creating a Reconcile Request")
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      crKey.Name,
					Namespace: crKey.Namespace,
				},
			}

			By("Directly invoking Reconciliation")
			// Invoke Reconcile
			_, err := reconciler.Reconcile(req)
			Expect(err).NotTo(HaveOccurred())

			By("Validating the details of the CRs deployment artifacts")
			// Validate the pod deployment
			pod := &corev1.Pod{}
			err = reconciler.Client.Get(testCtx, podKey, pod)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
