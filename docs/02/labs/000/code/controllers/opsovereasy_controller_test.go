package controllers

import (
	operatorsoverezv1alpha1 "github.com/mydomain/operators-over-ez/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"time"
)

var _ = Describe("CR Controller", func() {
	const timeout = time.Second * 60
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
			When("The specification timeout is set to a numeric value in seconds", func() {
				It("Should remain available for the specified timeout duration in seconds", func() {
					Expect(true).To(BeFalse())
				})
			})

			//SCENARIO: Log a user specified message before shutting down the busybox pod
			//GIVEN: An Operator instance
			//WHEN: the specification message is set to a string value
			//THEN: the busy box pod will log the message, from the message specification after the timeout duration has expired.
			When("The specification message is set to a string value", func() {
				It("Should log the message, from the message specification after the time out duration has expired", func() {
					Expect(true).To(BeFalse())
				})
			})

			//SCENARIO: Update status expired and logged when the busybox pod has expired
			//GIVEN: An Operator instance
			//WHEN: the busy box pod's duration has expired
			//THEN: set the expired status to true
			//AND: set the logged status to true
			When("The duration has expired", func() {
				It("Should set the expired status to true", func() {
					Expect(true).To(BeFalse())
				})

				It("Should set the logged status to true", func() {
					Expect(true).To(BeFalse())
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
					Expect(true).To(BeFalse())
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

			// Mock request to simulate Reconcile() being called on an event for a watched resource .
			By("Creating a Reconcile Request")
			req := reconcile.Request{
				NamespacedName: types.NamespacedName{
					Name:      crKey.Name,
					Namespace: crKey.Namespace,
				},
			}

			// Invoke Reconcile
			By("Directly invoking Reconciliation")
			_, err := opsOverEasyReconciler.Reconcile(req)
			Expect(err).NotTo(HaveOccurred())

			// Validate the pod deployment
			By("Validating the details of the CRs deployment artifacts")
			pod := &corev1.Pod{}
			err = opsOverEasyReconciler.Client.Get(testCtx, podKey, pod)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() corev1.PodPhase {
				pod := &corev1.Pod{}
				err = opsOverEasyReconciler.Client.Get(testCtx, podKey, pod)
				return pod.Status.Phase
			}, timeout, interval).Should(Equal(corev1.PodSucceeded))

			Expect(err).NotTo(HaveOccurred())
		})
	})
})
