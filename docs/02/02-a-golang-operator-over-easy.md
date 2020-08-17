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

## Environment Setup

Ensure lab pre-requisites have been met. See: [Lab Requirements](../01/03-lab-requirements.md)

## Lab Specifications

### Story (BDD/Gherkin Style)

**DESCRIPTION**

An Operator with a single busy box pod that logs a user specified message and shuts down after a user specified amount of time. If a duration or message are not specified, then both will be supplied by a REST API call. 

- **SCENARIO**: Shutdown the busybox pod after a user specified amount of time in seconds
  - **GIVEN**: A scaffolded operator
  - **AND**: an Operator instance
  - **WHEN**: the specification `timeout` is set to a numeric value in seconds
  - **THEN**: the busy box pod will remain available for the specified `timeout` in seconds,
  - **AND**: shutdown after the specified amount `timeout` duration

- **SCENARIO**: Log a user specified message before shutting down the busybox pod
  - **GIVEN**: A scaffolded operator
  - **AND**: an Operator instance
  - **WHEN**: the specification `message` is set to a string value
  - **THEN**: the busy box pod will log the message, from the `message` specification after the `timeout` duration has expired. 

- **SCENARIO**: Retrieve the `timeout` and `message` from a given REST API if one and/or the other is not supplied. 
  - **GIVEN**: A scaffolded operator
  - **AND**: an Operator instance
  - **WHEN**: the specification `message` OR `timeout` is NOT set
  - **THEN**: the busy box pod will supply these values from the following REST API: `GET http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response`

### Acceptance Criteria

- The CRD must have a `timeout` specification attribute
- The Operator instance must shut down after duration of `timeout` in seconds, has expired
- The CRD must have a `message` specification attribute
- The Operator instance must log the message `message` before the container has stopped
- The Operator instance must retrieve a `message` and `timeout` value from a REST API call (`GET http://my-json-server.typicode.com/keunlee/test-rest-repo/golang-lab00-response`), if both are not initially supplied on the Operator Instance. 

## Execution Strategy

In a nutshell, we want to start up a pod, running a busybox image for a specific duration and logging a user specific message. 

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

At this point, we've got a basic prototype of what we'd like the final deployment state of our Operator instance to be. In this case, it's a single busybox pod. 

The next step from here, is thinking about what our **domain specific operations** are. The previously generated pod YAML will not handle all of these operations as is. Rehashing requirements into domain specific operations: 

**If a message and duration are supplied, create a busybox pod with a duration and message** : This is pretty straightforward to automate. You just specify the `timeout` duration and `message` in the pods YAML. No real issues here. 

**If a message and duration are NOT supplied, then supply one from a REST API call, and then create a busybox pod with the duration and message**:  Since we've got a dynamic element at play here, we can automate this in code, w/in our Golang CR Controller. 

### II. Operator Scaffolding

Run the following to scaffold your operator and create a resource and controller. Say 'yes' to all prompts. 

```bash
# scaffold a new operator - over-ez-operator
operator-sdk init --domain=mydomain.com --repo=github.com/mydomain/over-ez-operator

# create new api and controller
operator-sdk create api --group=golang-op-lab00 --version=v1alpha1 --kind=Mycrd

# (you will be prompted the following) - create resource [y/n] y

# (you will be prompted the following) - create controller [y/n] y
```

One you run the above, you'll see a number of files generated. Of those files, the CR Implementation and controller: 

CR Implementation location: `api/v1alpha1/mycrd_types.go`

CR Controller location: `controllers/mycrd_controller.go`

### III. CR Definition Implementation

(1) Add specification attributes, per requirements. 

edit the file `api/v1alpha1/mycrd_types.go` by adding the `timeout` and `message` specifications w/in the `MycrdSpec struct` definition. It should look like the following: 

```golang
// MycrdSpec defines the desired state of Mycrd
type MycrdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

  // Accommodate requirements and acceptance criteria
	Timeout int32  `json:"timeout,omitempty"`
	Message string `json:"message,omitempty"`
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

You can validate specification updates on your CRD by examining the updated file: `crd/bases/golang-op-lab00.mydomain.com_mycrds.yaml`. You should see that your newly added specifications, `timeout` and `message`, have been added. 

### IV. TDD Setup

### V. CR Controller Implementation

### VI. Test Validation

```bash
go test ./controllers -timeout 30s -run ^TestAPIs$ -v
```

### VII. Deployment

