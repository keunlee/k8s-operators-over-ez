<!--
  - A Golang Operator - Over Easy
    - Overview
    - Lab Specification (BDD Style)
    - Step-by-Step Detailed Lab Walkthrough
      - Summary
      - Scaffolding
      - Writing Your Operator Specifications and Status
      - Writing Your Operator Controller Implementation
      - Unit Testing
      - End to End Testing
-->
# A Golang Operator - Over Easy

## Overview

This is a guided lab. In comparison to some of the other golang opeator labs that you will find in this guide, you will notice that this lab is heavier in annotations and explanations.

Please do not rush through this lab. You should leverage this lab to **learn**, get things to work, experiment and "break" stuff. It's okay to do that here. We will clean up the mess together!

That being said, these labs are meant to be run in a sandboxed development environment.

## Environment Setup

Ensure lab pre-requisites have been met. See: [Lab Requirements](../01/03-lab-requirements.md)

## Lab Specifications

### Story (BDD/Gherkin Style)

**TITLE**: Overeasy Operator Requirements

- **DESCRIPTION**
    - **AS A**: Developer
    - **I WANT**: An Operator with a single busybox pod that logs a user specified message and shuts down after a user specified amount of time. If a duration or message are not specified, then both will be supplied by a REST API call. 
    - **SO THAT**: I can demonstrate the encapsulation of operational knowlege, leveraging the Operator Design Pattern.  

- **SCENARIO 1**: Shutdown the busybox pod after a user specified amount of time in seconds
  - **GIVEN**: An Operator instance
  - **WHEN**: the specification `timeout` is set to a numeric value in seconds
  - **THEN**: the busy box pod will remain available for the specified `timeout` duration in seconds,

- **SCENARIO 2**: Log a user specified message before shutting down the busybox pod
  - **GIVEN**: An Operator instance
  - **WHEN**: the specification `message` is set to a string value
  - **THEN**: the busy box pod will log the message, from the `message` specification after the `timeout` duration has expired. 

- **SCENARIO 3**: Retrieve the `timeout` and `message` from a given REST API if one and/or the other is not supplied. 
  - **GIVEN**: An Operator instance
  - **WHEN**: the specification `message` OR `timeout` is NOT set
  - **THEN**: the busy box pod will supply these values from the following REST API: `GET http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response`

- **SCENARIO 4**: Update status `expired` and `logged` when the busybox pod has expired
  - **GIVEN**: An Operator instance
  - **WHEN**: the busy box pod's duration has expired
  - **AND**: the busy box pod has logged a message
  - **THEN**: set the operators `expired` status to `true`
  - **AND**: set the operators `logged` status to `true`

### Acceptance Criteria Notes

- The CRD must have a `timeout` specification attribute
- The Operator instance must shut down after the duration of `timeout` in seconds, has been reached
- The CRD must have a `message` specification attribute
- The Operator instance must log the message `message` before the container has stopped
- The CRD must have a `expired` status attribute
- The Operator must set the status of the busy box pod upon expiration, `expired`
- The CRD must have a `logged` status attribute
- The Operator must set the status of the busy box pod when logging a message, `logged`
- The Operator instance must retrieve a `message` and `timeout` value from a REST API call (`GET http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response`), if both are not initially supplied on the Operator Instance. 

## Execution Strategy

In a nutshell, we want to start up a pod, running a busybox image for a specific duration and logging a user specific message, and then setting our Operator's status. 

We'll want our Operator to provision our pod with the necessary attribute specifications, eventually. 

For now, our strategy to reach the end state is detailed as followed: 

- **I - Prototyping** - Create a YAML specification for a pod which runs for a specified amount of time and logs a specific message. Do this to validate our design. We'll eventually want our Operator controller implementation to dynamically set the pods timeout duration and log message. For now, we will validate our prototype. 

- **II - Operator Scaffolding** - Scaffold a Golang Operator to give us an initial template for our CRD and Resource Controller

- **III - CR Definition Implementation** - Add the `timeout` attribute to our CRD.

- **IV - TDD Setup** - Create a Unit Test file for our Controller to validate our requirements leveraging TDD (Test Driven Design). We will validate the tests as we implement our controller. 

- **V - CR Controller Implementation**- Implement our Resource Controller logic to help fulfill the Story and Acceptance Criteria.

- **VI - Test Validation** - Validate our Unit Tests. Sanity check our Operator so that it is indeed operating as intended. 

- **VII - Deployment** - Deploy the Operator to your Kubernetes cluster

> :information_source: CR is an acronym for "Custom Resource"

### I. Prototyping

(1) Build a Proof of Concept

Let's begin by creating a project namespace in our cluster. 

```bash
kubectl create ns golang-op-lab-00
```

set the current context to newly created namespace

```bash
kubens golang-op-lab-00
```

Create the yaml for a pod which will start a busybox container and run for a specified duration, 15 seconds, and logs the message "hello world".

```bash
# create the pod yaml
kubectl run busybox --image=busybox --restart=Never --dry-run -o yaml -- /bin/sh -c 'sleep 15; echo "hello world"' > golang-op-lab-00-pod.yam
```

Running the following will yield the following generated yaml contents: 

```yaml
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    run: busybox
  name: busybox
spec:
  containers:
  - args:
    - /bin/sh
    - -c
    - sleep 15; echo "hello world"
    image: busybox
    name: busybox
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Never
status: {}
```

To deploy the pod, watch it's change in status after the set duration, and view the pods logs:  

```bash
# deploy the pod
kubectl apply -f golang-op-lab-00-pod.yaml

# watch for changes on the pod, ctrl-c to exit the watch
watch kubectl get po

# display log messages
kubectl logs busybox -c busybox
```

(2) Identify Domain Specific Operations

At this point, we've got a basic prototype of what we'd like the final deployment state of our Operator in tried to
The next step from here, is thinking about what our **domain specific operations** are. The previously generated pod YAML will not handle all of these operations as is. Rehashing requirements into domain specific operations: 

- **If a message and duration are supplied, create a busybox pod with a duration and message** : This is pretty straightforward to automate. You just specify the `timeout` duration and `message` in the pods YAML. No real issues here. 

- **If a message and duration are NOT supplied, then supply one from a REST API call, and then create a busybox pod with the duration and message**:  Since we've got a dynamic element at play here, we can automate this in code, w/in our CR Controller. 

We now have a purpose for building an operator, which will automate our domain specification operations. 

### II. Operator Scaffolding

To begin codifying our domain specific operations, run the following to scaffold your operator and create a resource and controller. Say 'yes' to all prompts. 

```bash
# scaffold a new operator - over-ez-operator
operator-sdk init --domain=mydomain.com --repo=github.com/mydomain/operators-over-ez

# create new api and controller
 operator-sdk create api --group=operators-over-ez --version=v1alpha1 --kind=OpsOverEasy

# (you will be prompted the following) - create resource [y/n] y

# (you will be prompted the following) - create controller [y/n] y
```

One you run the above, you'll see a number of files generated. Of those files, the CR Implementation and controller: 

- CR Implementation location: `api/v1alpha1/opsovereasy_types.go`

- CR Controller location: `controllers/opsovereasy_controller.go`

These resources will be particularly important going forward. 

### III. CR Definition Implementation

(1) Add specification and status attributes, per requirements. 

In this step, you will define your operator specifications and status attributes. 

**Operator Specifications**:

These are user defined properties of your operator. When you add a specification, you typically implement your operator to respond/react off of the specification that is set. 

<ins>Example 1</ins>

*We set the specifications `timeout` and `message` to finite/non-null values in our operator instance yaml. When we create the operator instance, our operator is expected to then deploy a busybox pod, which will run the busybox container for the `timeout` duration and log the `message`.* 

<ins>Example 2</ins>

*We leave the specifications for `timeout` and `message` unset in our operator instance yaml. When we create the operator instance, our operator is expected to then obtain values for `timeout` and `message` from a REST API call and then deploy a busybox pod, which will run the busybox container for the `timeout` duration and log the `message`.* 

<ins>Do This</ins>

Edit the file `api/v1alpha1/opsovereasy_types.go` by adding the `timeout` and `message` specifications w/in the `OpsOverEasySpec struct` definition. It should look like the following: 

```golang
// OpsOverEasySpec defines the desired state of Mycrd
type OpsOverEasySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Accommodate requirements and acceptance criteria

	// Timeout duration, in seconds
	Timeout int32 `json:"timeout,omitempty"`

	// Log message to output before container expires
	Message string `json:"message,omitempty"`
}
```

Our **CR Controller**, must implement the logic, workflow and functionality to facilitate these examples, which are based on our requirements. 

**Operator Statuses**:

These are user defined properties for checking the status/state of your operator deployment. When you add a status, you typically implement your operator to provide the state/value of the status. 

<ins>Example</ins>

*When we deploy our operator, we expect a busybox pod to be spun up, log message, and then spin down after set duration. When the pod has reached this state, then the operator will need to set it's status for `expired` and `logged` to `true`.* 

<ins>Do This</ins>

Edit the file `api/v1alpha1/opsovereasy_types.go` by adding the `expired` and `logged` statuses w/in the `OpsOverEasyStatus struct` definition. It should look like the following: 

```golang
// MycrdStatus defines the observed state of Mycrd
type OpsOverEasyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Denotes if the time duration has expired
	TimeoutExpired bool `json:"expired,omitempty"`

	// Denotes if a message was logged
	MessageLogged bool `json:"logged,omitempty"`
}
```

(2) Regenerate/update resource and manifest

Always do this after you update your custom resource in `*_types.go`

```bash
# generate/update code for resource types
make generate

# generate/update manifests for the CRD
make manifests
```

You can validate specification updates on your CRD by examining the updated file: `config/crd/bases/operators-over-ez.mydomain.com_opsovereasies.yaml`. You should see that your newly added specifications, `timeout` and `message`, have been added. 

### IV. TDD Setup

For this lab we will be taking a TDD (Test Driven Design) approach to building and implementing our operator. 

(1) Update the Test Suite

When we generated our custom resource definition and controller, we also generated a stub for a test suite: `controllers\suite_test.go`

<ins>Do This</ins>

We are going to want update our Test Suite file to look like the following:  https://bit.ly/2Qi9ZRQ

When observing the additions made to the generated stub: 

- Added imports to additional packages that we will be leveraging ![Screenshot from 2020-08-24 15-26-01](https://user-images.githubusercontent.com/61749/91093790-87816400-e61f-11ea-91a5-b9ca1cf9e6fe.png)

- Added additional variables and helper functions that we'd like to make available to our controller tests. ![Screenshot from 2020-08-24 15-26-35](https://user-images.githubusercontent.com/61749/91093791-8819fa80-e61f-11ea-8e8f-d2fe926d2524.png)

- Since we'll need to validate the internals of our deployment (i.e. such as pod phase, etc.), we'll be want to run these tests against a live cluster. To automate a test cluster, one suggestion would be to automate the provisioning of a cluster, leveraging a lightweight cluster implementation such as [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/) or [K3S/K3D](https://k3d.io/) ![Screenshot from 2020-08-24 22-14-59](https://user-images.githubusercontent.com/61749/91118825-7ce2c100-e657-11ea-83e3-6b6840540bdc.png)

- An instance of a controller reconciler we'll be using in our test ![Screenshot from 2020-08-24 15-27-11](https://user-images.githubusercontent.com/61749/91093794-88b29100-e61f-11ea-8df5-f3e855770e57.png)

- We can automate the uninstall of a CRD after it has been tested on our cluster. ![Screenshot from 2020-08-24 22-15-36](https://user-images.githubusercontent.com/61749/91118827-7d7b5780-e657-11ea-9c4c-2f248ac294e5.png)

(2) Add Test Controller and Test Stubs

<ins>Do This</ins>

Create the following file: `controllers/opsovereasy_controller_test.go`

Copy the following contents of test implementation stub to the file:

```golang
package controllers

import (
	operatorsoverezv1alpha1 "github.com/mydomain/operators-over-ez/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("CR Controller", func() {

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
					Expect(true).To(BeFalse())
				})

				It("Should shutdown after the specified amount timeout duration", func() {
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
})
```

Notice how we we've created stubs for each corresponding BDD scenario. 

As we implement each BDD scenario we will also leverage these stubs to test and validate our scenario requirements. 

Currently when these tests are executed they will fail as expected. 

This will change as we implement our controller and update these tests. 

> You can find a full implementation of this file here: https://bit.ly/2YxaAnb

(3) Run Tests

You can run these tests by running the following at the terminal of the root directory of the lab: 

```bash
make test
```

All tests should fail at this point. This is as  expected. 

### V. CR Controller Implementation

### VI. Test Validation

```bash
go test ./controllers -timeout 30s -run ^TestAPIs$ -v
```

### VII. Deployment

