package controllers

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var pod = &corev1.Pod{}
var currUuid string

func createReconcileRequest() error {
	// make request to Reconcile
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      crdInstance.Name,
			Namespace: crdInstance.Namespace,
		},
	}

	podKey := types.NamespacedName{
		Namespace: "default",
		Name:      crdInstance.Name + "-pod",
	}

	// Invoke Reconcile
	_, err := opsOverEasyReconciler.Reconcile(req)
	if err != nil {
		return err
	}

	// Validate the pod deployment by retrieving it
	pod = &corev1.Pod{}
	err = opsOverEasyReconciler.Client.Get(testCtx, podKey, pod)
	return err
}

func getPodLogs(pod corev1.Pod) string {
	podLogOpts := corev1.PodLogOptions{}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return "error in getting access to K8S"
	}
	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(testCtx)
	if err != nil {
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}

var _ = Describe("CR Controller", func() {
	Context("BDD Test Scenarios", func() {
		Context("CR Instance with Specifications Provided", func() {

			BeforeEach(func() {
				currUuid = string(uuid.NewUUID())
				crdInstance = getCrd(true, currUuid)

				Eventually(func() error {
					err := k8sClient.Create(testCtx, crdInstance)
					return err
				}, timeout, interval).Should(Succeed())
			})

			AfterEach(func() {
				Eventually(func() error {
					err := k8sClient.Delete(testCtx, crdInstance)
					return err
				}, timeout, interval).Should(Succeed())
			})

			//SCENARIO 1: Shutdown the busybox pod after a user specified amount of time in seconds
			//GIVEN: An Operator instance
			//WHEN: the specification `timeout` is set to a numeric value in seconds
			//THEN: the busy box pod will remain available for the specified `timeout` duration in seconds,
			When("The specification `timeout` is set to a numeric value in seconds", func() {
				It("Should remain available for the specified timeout duration in seconds", func() {
					Expect(crdInstance.Spec.Timeout).Should(Equal(int32(podDuration)))

					err := createReconcileRequest()
					Expect(err).NotTo(HaveOccurred())

					Eventually(func() corev1.PodPhase {
						pod = &corev1.Pod{}
						podKey := getPodKey(currUuid)
						err = opsOverEasyReconciler.Client.Get(testCtx, podKey, pod)
						return pod.Status.Phase
					}, timeout, interval).Should(Equal(corev1.PodSucceeded))
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			//SCENARIO 2: Log a user specified message before shutting down the busybox pod
			//GIVEN: An Operator instance
			//WHEN: the specification message is set to a string value
			//THEN: the busy box pod will log the message, from the message specification after the timeout duration has expired.
			When("The specification message is set to a string value", func() {
				It("Should log the message, from the message specification after the time out duration has expired", func() {
					Expect(crdInstance.Spec.Message).Should(Equal("message"))

					err := createReconcileRequest()
					Expect(err).NotTo(HaveOccurred())

					Eventually(func() corev1.PodPhase {
						pod = &corev1.Pod{}
						podKey := getPodKey(currUuid)
						err = opsOverEasyReconciler.Client.Get(testCtx, podKey, pod)
						return pod.Status.Phase
					}, timeout, interval).Should(Equal(corev1.PodSucceeded))
					Expect(err).NotTo(HaveOccurred())

					logs := getPodLogs(*pod)
					Expect(logs).NotTo(BeEmpty())
					Expect(logs).To(ContainSubstring("message"))
					Expect(err).NotTo(HaveOccurred())
				})
			})

			//SCENARIO 4: Update status expired and logged when the busybox pod has expired
			//GIVEN: An Operator instance
			//WHEN: the busy box pod's duration has expired
			//THEN: set the expired status to true
			//AND: set the logged status to true
			When("The duration has expired", func() {
				It("Should set the expired and logged status to true", func() {
					Expect(crdInstance.Status.TimeoutExpired).Should(BeFalse())
					Expect(crdInstance.Status.MessageLogged).Should(BeFalse())

					err := createReconcileRequest()
					Expect(err).NotTo(HaveOccurred())

					Eventually(func() corev1.PodPhase {
						pod = &corev1.Pod{}
						podKey := getPodKey(currUuid)
						err = opsOverEasyReconciler.Client.Get(testCtx, podKey, pod)
						return pod.Status.Phase
					}, timeout, interval).Should(Equal(corev1.PodSucceeded))
					Expect(err).NotTo(HaveOccurred())

					err = createReconcileRequest()
					Expect(err).NotTo(HaveOccurred())

					crdInstance = getCrd(true, currUuid)
					crKey := getCrKey(currUuid)

					err = k8sClient.Get(testCtx, crKey, crdInstance)
					Expect(err).NotTo(HaveOccurred())

					Expect(crdInstance.Status.TimeoutExpired).Should(BeTrue())
					Expect(crdInstance.Status.MessageLogged).Should(BeTrue())
				})
			})
		})

		Context("CR Instance with no Specifications Provided", func() {
			BeforeEach(func() {
				currUuid = string(uuid.NewUUID())
				crdInstance = getCrd(false, currUuid)

				Eventually(func() error {
					err := k8sClient.Create(testCtx, crdInstance)
					return err
				}, timeout, interval).Should(Succeed())
			})

			AfterEach(func() {
				Eventually(func() error {
					err := k8sClient.Delete(testCtx, crdInstance)
					return err
				}, timeout, interval).Should(Succeed())
			})

			//SCENARIO 3: Retrieve the timeout and message from a given REST API if one and/or the other is not supplied.
			//GIVEN: An Operator instance
			//WHEN: the specification message OR timeout is NOT set
			//THEN: the busy box pod will supply these values from the following REST API: GET http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response
			When("The specification message OR timeout is NOT set", func() {
				It("Should supply these values from the following REST API: GET http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response", func() {
					Expect(crdInstance.Spec.Timeout).Should(Equal(int32(0)))
					Expect(crdInstance.Spec.Message).Should(Equal(""))

					err := createReconcileRequest()
					Expect(err).NotTo(HaveOccurred())

					Eventually(func() corev1.PodPhase {
						pod = &corev1.Pod{}
						podKey := getPodKey(currUuid)
						err = opsOverEasyReconciler.Client.Get(testCtx, podKey, pod)
						return pod.Status.Phase
					}, timeout, interval).Should(Equal(corev1.PodSucceeded))
					Expect(err).NotTo(HaveOccurred())

					crdInstance = getCrd(false, currUuid)
					crKey := getCrKey(currUuid)
					err = k8sClient.Get(testCtx, crKey, crdInstance)
					Expect(err).NotTo(HaveOccurred())

					Expect(crdInstance.Spec.Timeout).Should(Equal(int32(5)))
					Expect(crdInstance.Spec.Message).Should(Equal("domain specific operational knowledge is king"))
				})
			})
		})
	})
})
