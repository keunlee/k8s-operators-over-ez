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

### Stories (BDD/Gherkin Style)

**DESCRIPTION**

An Operator with a single busy box pod that logs a user specified message and shuts down after a user specified amount of time

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

### Acceptance Criteria

- the Operator must have a `timeout` specification attribute
- the Operator instance must shut down after duration of `timeout` in seconds, has expired
- the Operator must have a `message` specification attribute
- the Operator instance must log the message `message` before the container has stopped

## Execution Strategy

In a nutshell, we want to start up a pod, running a busybox image for a specific duration and logging a user specific message. 

We'll want our Operator to provision our pod with our user specified attribute selections, eventually. 

For now, our strategy to reach the end state is detailed as followed: 

- **I - Prototyping** - Create a YAML specification for a pod which runs for a specified amount of time and logs a specific message. Do this to validate our design. We'll eventually want our Operator controller implementation to dynamically set the pods timeout duration and log message. For now, we will validate our prototype. 

- **II - Operator Scaffolding** - Scaffold a Golang Operator to give us an initial template for our CRD and Resource Controller

- **III - TDD Setup** - Create a Unit Test file for our Controller to validate our requirements leveraging TDD (Test Driven Design). We will validate the tests as we implement our controller. 

- **IV - CR Definition Implementation** - Add the `timeout` attribute to our CRD.

- **V - CR Controller Implementation**- Implement our Resource Controller logic to help fulfill the Story and Acceptance Criteria.

- **VI - Test Validation** - Validate our Unit Tests. Sanity check our Operator so that it is indeed operating as intended. 

- **VII - Deployment** - Deploy the Operator to your Kubernetes cluster

> :information_source: CR is an acronym for "Custom Resource"

## I. Prototyping

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

If we deploy this yaml, we'll see that it will run for 15 seconds, log a message to the console and shutdown afterwards. 

To deploy the pod, watch it's change in status after the set duration, and view the pods logs:  

```bash
# deploy the pod
kubectl apply -f golang-op-lab-00-pod.yaml

# watch for changes on the pod, ctrl-c to
watch kubectl get po

# display log messages
kubectl logs busybox -c busybox
```

## II. Operator Scaffolding

## III. TDD Setup

## IV. CR Definition Implementation

## V. CR Controller Implementation

## VI. Test Validation

## VII. Deployment

# Going Forward

