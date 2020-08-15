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

### Story

- **DESCRIPTION**: An Operator with a single busy box pod that shuts down after a user specified amount of time
  - **GIVEN**: A scaffolded operator
  - **WHEN**: the specification `timeout` is added as an attribute to the operator
  - **AND**: the specification `timeout` is set to a numeric value in seconds
  - **AND**: and an Operator instance is created
  - **THEN**: the busy box pod will remain available for the specified `timeout` in seconds,
  - **AND**: log the message, `busybox pod expired` upon `timeout` expiration 
  - **AND**: shutdown after the specified amount `timeout` duration

### Acceptance Criteria

- the Operator must have a `timeout` specification attribute
- the Operator instance must start a busybox pod for the duration of `timeout` in seconds
- the Operator instance must shut down after duration of `timeout` in seconds, has expired
- the Operator instance must log a message before expiration

## Execution Strategy

In a nutshell, we want to start up a pod, running a busybox image for a specific duration. But we want our Operator to do this for us, eventually. Our strategy to reach the end state is detailed as followed: 

- **I - Design** - Create a YAML specification for a pod which runs for a specified amount of time. Do this to validate our design and to validate that our busybox pod can run for a set duration. 

- **II - Scaffolding** - Scaffold a Golang Operator to give us an initial template for our CRD and Resource Controller

- **III - TDD** - Create a Unit Test file for our Controller to validate our requirements leveraging TDD (Test Driven Design). We will validate the tests as we implement our controller. 

- **IV - CR Definition Implementation** - Add the `timeout` attribute to our CRD.

- **V - CR Controller Implementation**- Implement our Resource Controller logic to help fulfill the Story and Acceptance Criteria.

- **VI - Test Validation** - Validate our Unit Tests. Sanity check our Operator so that it is indeed operating as intended. 

- **VII - Deployment** - Deploy the Operator to your Kubernetes cluster

> :information_source: CR is an acronym for "Custom Resource"

## I. Design

Let's begin by creating a project namespace in our cluster. 

```bash
kubectl create ns golang-op-lab-00
```

set the current context to newly created namespace

```bash
kubens golang-op-lab-00
```

Let's try to create the yaml for a pod which will start a busybox container and run for a specified duration, 15 seconds.  

```bash
# create the pod yaml
kubectl run busybox --image=busybox --restart=Never --dry-run -o yaml -- /bin/sh -c 'sleep 15' > golang-op-lab-00-pod.yaml
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
    - sleep 10
    image: busybox
    name: busybox
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Never
status: {}
```

If we deploy this yaml, we'll see that it will run for 15 seconds and shutdown afterwards. To deploy the pod and watch it's change in status after the set duration:  

```bash
# deploy the pod
kubectl apply -f golang-op-lab-00-pod.yaml

# watch for changes on the pod, ctrl-c to
watch kubectl get po
```


## II. Scaffolding

## III. TDD

## IV. CR Definition Implementation

## V. CR Controller Implementation

## VI. Test Validation

## VII. Deployment

## Going Forward

