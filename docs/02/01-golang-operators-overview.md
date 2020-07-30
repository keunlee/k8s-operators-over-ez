# Golang Operators

## Prerequisites

- Review the [introduction](../01/01-introduction.md), if you have not already
- Familiarity with Golang. If you are not familiar with Go, atleast go through this resource: https://tour.golang.org/welcome/1

## Agenda

This section will cover the following: 

- Review options for working with Operators in Golang
- Useful Golang API References for working with CRDs, Controllers, and the [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- Review the [Reconciliation Cycle](../01/01-introduction.md#how-do-operators-work) with Golang semantics

Afterwards, you will take the plunge in a guided walkthrough. 

# Golang Operator Resources

## Development Libraries

Two resources you can check out for writing Operators in Golang: 

- [Operator Framework](https://operatorframework.io/)
- [Kubebuilder](https://book.kubebuilder.io/quick-start.html)

> In this section, as we discuss Golang Operators, we will be referring to the __Operator SDK__ libary, as we continue to discuss. 

In case you are curious of some of the differences between the two, here's a recap: [What is the difference between kubebuilder and operator-sdk?](https://github.com/operator-framework/operator-sdk/issues/1758#issuecomment-517432349)

## Useful Golang API References

>NOTE: By no means, I'm not expecting you to memorize or review in detail what the APIs below are and do. I expect you to just be aware of the location of these API docs, as they may come in handy as you build out your own operator(s) and would like to understand better what these APIs are, as the Golang Operator libaries make heavy use of them, without a ton of explanation, rhyme, or reason of what they are and what they do. 

The APIs below are commonly leveraged by Kubernetes. Subsequently, they are also commonly leveraged by Golang Operator development libraries. You will usage of these libaries, as well as others, when scaffolding your Operators for development. Many of the API reference docs come with inline examples of how to use them in code. This is just a small subset of what's actually available out here in the community. 

- Controller Runtime
    - https://pkg.go.dev/github.com/kubernetes-sigs/controller-runtime?tab=doc
        - The Kubernetes controller-runtime Project is a set of go libraries for building Controllers. It is leveraged by Kubebuilder and Operator SDK. Both are a great place to start for new projects. See Kubebuilder's Quick Start to see how it can be used. 
        - Package controllerruntime alias' common functions and types to improve discoverability and reduce the number of imports for simple Controllers.
        - see: https://github.com/kubernetes-sigs/controller-runtime

    - https://pkg.go.dev/sigs.k8s.io/controller-runtime?tab=doc
        - Package controllerruntime provides tools to construct Kubernetes-style controllers that manipulate both Kubernetes CRDs and aggregated/built-in Kubernetes APIs. It defines easy helpers for the common use cases when building CRDs, built on top of customizable layers of abstraction. Common cases should be easy, and uncommon cases should be possible. In general, controller-runtime tries to guide users towards Kubernetes controller best-practices.
  
- https://pkg.go.dev/k8s.io/apimachinery?tab=overview
  - Scheme, typing, encoding, decoding, and conversion packages for Kubernetes and Kubernetes-like API objects.

- https://pkg.go.dev/k8s.io/api?tab=overview
  - Schema of the external API types that are served by the Kubernetes API server.
  - Most notably, you will be able to find the types for constructing resources such as (to name a few):
    - Pods
    - Services
    - Deployments
    - Replicasets
    - Daemonsets

# The Reconciliation Cycle - Revisited

> In the introduction, we presented the reconciliation cycle in an operator controller as followed: 

![](../assets/resource-controller-reconciliation-cycle.png)

> In this section, we will recap the Reconciliation Cycle in more detail, with respect to a Golang Operator Controller. The embellishments to the Reconciliation Cycle are shown below: 

![](../assets/resource-controller-reconciliation-cycle-golang-operators.png)


## Observe/Watch

In this phase, the controller observes the state of th cluster. Typically this is initiated by observing the events on the custom resource instance. These events are usually subscribed from the custom resource controller. Consider this to be similar in ways to a pub/sub mechanism between the resource controller and cluster. 

TODO: add more

## Analyze

In this phase, the resource controller compares the current state of the resource instance to the desired state. The desired state is typically reflective of what is specified in the `spec` attributes of the resource. 

TODO: add more

## Act/Reconcile

In this phase, the resource controller performs all necessary actions to make the current resource state match the desired state. This is called reconciliation, and is typically where operational knowledge is implemented (i.e. business/domain logic).

TODO: add more