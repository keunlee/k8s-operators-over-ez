<!--
- Introduction
  - How to Learn to Write Operators Using this Resource
  - What's an Operator?
  - Types of Operators
  - Operator Tooling & Resources  
-->

# Introduction

This guide serves all but one purpose: 

**To help assist people like you to understand how to write kubernetes operators**

## How to Learn to Write Operators Using this Resource

The approach to learning from this resource is composed into the following components (not in any particular order). 

- direct and straighforward explanation of things
- illustrations
- links to other related resources
- step by step instructions
- labs
- repetition

The views and opinions w/in this guide are that of the author. You may not agree with all the views expressed. Feel free to raise a defect/issue/etc whereever you encounter one. 

You will see that the labs are pretty basic in regards to a functional aspect of an operator. While the complexity of the labs may increase, the author's hope is that it does in a comprehensible and easy to understand fashion.  

The point is to illustrate, explain, and eventually, have you get to the point of not thinking about language and api semantics, but to get you to focus on what you want to build and what you want your operator to do. 

## What's an Operator?

From the author's point of view, an operator allows you to encapsulate a set/grouping of kubernetes deployable artifacts (i.e. pods, deployments, daemonsets, replicasets, services, configmaps, etc), by creating a [CRD/Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) and a custom [Controller](https://kubernetes.io/docs/concepts/architecture/controller/) implementation which may drive deployment and/or business logic of your CRD deployment instance(s).

The advantages of an operator can be seen when you start thinking about how you manage your applications deployed into a Kubernetes cluster. Typically, you manage the deployment of each individual artifact (i.e. pods, deployments, daemonsets, replicasets, services, configmaps, etc) which may compose the entirety of the application you deploy to kubernetes cluster.

With operators, these artifacts are packaged up, and only expose the necessities (i.e. configurations, specifications, etc. ) of each artifact in one place of configuration -- your "Operator" instance. 

![](../assets/conventional-vs-operators.png)

### Resources

- https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
- https://coreos.com/blog/introducing-operators.html
- https://www.openshift.com/blog/operator-framework-moves-to-cncf-for-incubation
- https://www.openshift.com/blog/kubernetes-operators-best-practices


## Types of Operators

With regards to the Red Hat's Operator SDK, the following Operator Types are supported: 

- Golang Operators
- Helm Operators
- Ansible Operators

But wait, there's more! Outside of the Operator SDK ecosystem: 

- Java Operators
- Python Operators
- Javascript/Typescript Operators
- JSONNET Operators

It should be noted there is a certain amount of flexibility and choice that you have

## Operator Tooling & Resources

### Frameworks to Help you Create Operators (Not an Exhaustive List)

- https://kudo.dev/
- https://book.kubebuilder.io/
- https://metacontroller.app/
- https://github.com/operator-framework/getting-started
- https://github.com/fabric8io/kubernetes-client

For the entirety of this documentation, the author has decided to leverage Coreos's [Operator Framework](https://github.com/operator-framework). 