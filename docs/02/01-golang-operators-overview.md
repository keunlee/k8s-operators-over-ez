# Golang Operators

## Prereqisites

- Review the [introduction](../01/01-introduction.md), if you have not already
- Familiarity with Golang. If you are not familiar with Go, atleast go through this resource: https://tour.golang.org/welcome/1

## Agenda

- Review options for working with Operators in Golang
- Additional Golang References for working with CRDs, Controllers, and the [Controller Runtime](https://github.com/kubernetes-sigs/controller-runtime)
- Review the [Reconciliation Cycle](../01/01-introduction.md#how-do-operators-work) with Golang semantics

Afterwards, you will take the plunge in a guided walkthrough. 

# Golang Operator Resources

Two resources you can check out for writing Operators in Golang: 

- [Operator Framework](https://operatorframework.io/)
- [Kubebuilder](https://book.kubebuilder.io/quick-start.html)

In this documentation, you will be leveraging the Operator Framework. 

In case you are curious of some of the differences between the two, here's a recap: [What is the difference between kubebuilder and operator-sdk?](https://github.com/operator-framework/operator-sdk/issues/1758#issuecomment-517432349)


# The Reconciliation Cycle - Revisited

![](../assets/resource-controller-reconciliation-cycle-golang-operators.png)

Let's recap the Reconciliation Cycle in a bit more detail. 

## Observe/Watch

In this phase, the controller observes the state of th cluster. Typically this is initiated by observing the events on the custom resource instance. These events are usually subscribed from the custom resource controller. Consider this to be similar in ways to a pub/sub mechanism between the resource controller and cluster. 

## Analyze

In this phase, the resource controller compares the current state of the resource instance to the desired state. The desired state is typically reflective of what is specified in the `spec` attributes of the resource. 

## Act/Reconcile

In this phase, the resource controller performs all necessary actions to make the current resource state match the desired state. This is called reconciliation, and is typically where operational knowledge is implemented (i.e. business/domain logic).