# I. Operator Development Lab - Tasks

- [x] deploy a busy box pod
- [x] deploy multiple busy box pods
- [ ] deploy a deployment
- [ ] deploy an nginx deployment + service
- [ ] deploy a custom web app and make it accessible

- [ ] helm based operators

- [ ] ansible based operators

# II. Basics

## Install Operator Framework

go here: https://sdk.operatorframework.io/docs/install-operator-sdk/

# III. Create A Single-Pod Busy Box Operator

The most basic operator that you can build using the operator sdk. For the most part, it is essentially out of the box. 

By default, the operator sdk scaffolds will give you the following artifacts upon running basic scaffolding commands (i.e. create crd, add controller, etc.)

- CRD (custom resource definition
  - yaml file - CRD yaml file
  - golang file - this drives the code generation of the aformentioned CRD yaml file
- a Controller for your resource -- as a golang file

An operator instance deployment will do the following: 

- deploy one busy box pod

This is all done without writing additional code 

## Create Operator 

```bash
# create and select namespace
kubectl create ns task001
kubens task001

# create project
operator-sdk new task-001-operator --repo=github.com/keunlee/task-001-operator
cd task-001-operator

# create crd
operator-sdk add api --api-version=task-001.thekeunster.local/v1alpha1 --kind=Task001
operator-sdk generate k8s
operator-sdk generate crds

# add controller
operator-sdk add controller --api-version=task-001.thekeunster.local/v1alpha1 --kind=Task001

# deploy operator
operator-sdk apply -f deploy/crds/task-001.thekeunster.local_task001s_crd.yaml

# run locally outside of cluster
# make sure you are in your operator directory
operator-sdk run local --watch-namespace=default
```

## Deploy Operator Instance

```bash
kubectl create -f deploy/crds/task-001.thekeunster.local_v1alpha1_task001_cr.yaml
```

## Validate Operator Instance

### Validate Pod(s)
```bash
kubectl get po
```

yields: 

```bash
NAME                  READY   STATUS    RESTARTS   AGE
example-task001-pod   1/1     Running   0          46s
```

### Validate the Operator Instance

```bash
kubectl get Task001
```

yields: 

```bash
NAME              AGE
example-task001   19m
```

## Testing

### Unit Testing

### E2E Testing

# IV. Create A Multi-Pod Busy Box Operator 

## Operator Requirements

```
(1) 

DESCRIPTION: Add an operator specification, "NumberOfPods", which is used to specify the number of busybox pods to create. 
GIVEN: The operator specification, "NumberOfPods"
WHEN: "NumberOfPods" is set to 'n', where n > 0
THEN: The operator will scale the number of busybox pods up/down to n

(2) 

DESCRIPTION: Add an operator status, "ListedPods", which is used to track the names of all deployed pods by the operator. 
GIVEN: The operator status, "ListedPods"
WHEN: an operator instance is deployed with n pods (where n > 0)
THEN: "ListedPods" will be assigned a named string based list of all deployed pods by the operator. 
```

## Procedure

### 1. Generate Boilerplate

The custom script will run boiler plate operations to create an operator controller and types

```bash
OPERATOR_NAME=task-002 CRD_NAME=Task002 source automation/create-golang-operator.sh
```
### 2. Add Type Specs and Status According to Requirements

(1) open the file `pks/apis/task002/v1alpha1/task002_types.go`

(2) modify the `Task002Spec` and `Task002Status` structs respectively in the following fashion. 

```go
// Task002Spec defines the desired state of Task002
type Task002Spec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	NumberOfPods	int 	`json:"numpods,omitempty"`
}

// Task002Status defines the observed state of Task002
type Task002Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ListedPods	[]string	`json:"listedpods,omitempty"`
}
```

This is where you add your operator specifications and status attributes. 

**Operator Specifications**: 

These are user defined properties of your operator. When you add a specification, you typically implement your operator to respond/react off of the specification that is set. So for example: We set the specification `'numpods` in our operator instance `yaml` to the value of 3. When we create the operator instance, our operator, according to our requirements, is expected to then deploy 3 busybox pods. Out controller logic, must imlement the functionality to facilitate to this. 

**Operator Statuses**: 

These are user defined properties for checking the status/state of your operator deployment. When you add a status, you typically implement your operator to provide the state/value of the status. For example: when we deploy the operator with a specification `numpods` equal to 3, we expect that there will be three pods that are spun up, all with unique names. The status `listedpods` will need to list the names of the pods currently deployed. Our controller logic must implement the functionality required to fullfill the values of `listedpods`

(3) re-generate your CRD files to accomodate the new attributes:

```bash
operator-sdk generate k8s
operator-sdk generate crds
```

Once you've regenerated your CRD file, you must then deploy your operator. 

```bash
# build and deploy the operator
kubectl create -f deploy/crds/task-002.thekeunster.local_task002s_crd.yaml
```

### 3. Add/Update Controller Logic According to Requirements

open the controller file: `pkg/controller/task002/task002_controller.go`

locate the _Reconcile_ function and modify it with the following skeleton

```go
func (r *ReconcileTask002) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Task002")

	// Fetch the Task002 instance
	instance := &task002v1alpha1.Task002{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

    // (1) Retrieve a list of existing pods in our namespace

    // (2) Count the pods that are pending or running as available

    // (3) Update the status if necessary

    // (4) scale pods down to 'NumberOfPods' specification

    // (5) scale pods up to 'NumberOfPods' specification
    
	return reconcile.Result{}, nil
}
```

**(0)**





# V. Tips/Tricks

- setup an IDE for development. options: 
  - GoLand: https://medium.com/@auscunningham/debug-kubernetes-operator-sdk-locally-in-goland-27b7909c417a
  - VS Code: https://medium.com/@auscunningham/debug-kubernetes-operator-sdk-locally-using-vscode-a233aa7c750e
- TODO: show how validate your IDE setup in GoLand and VS Code
- Learn how to interact with your cluster via a controller through the controller-runtime resources: 
  - https://github.com/kubernetes-sigs/controller-runtime
  - https://godoc.org/github.com/kubernetes-sigs/controller-runtime/pkg

# VI. Anatomy of a Controller

## Watching: `func add(mgr manager.Manager, r reconcile.Reconciler)`

## Reconciling: `func (r *ReconcileTask001) Reconcile(request reconcile.Request)`

## Resource Creation

# VII. Appendix

## Boiler Plate

```bash
# create project
operator-sdk new memcached-operator --repo=github.com/example-inc/memcached-operator
cd memcached-operator

# create crd
operator-sdk add api --api-version=cache.example.com/v1alpha1 --kind=Memcached
operator-sdk generate k8s
operator-sdk generate crds

# add controller
operator-sdk add controller --api-version=cache.example.com/v1alpha1 --kind=Memcached

# build and run operator
kubectl create -f deploy/crds/cache.example.com_memcacheds_crd.yaml
```

## Run Operator Instance

```bash
# run locally outside of cluster
# make sure you are in your operator directory
operator-sdk run local --watch-namespace=default

# create an instance of the operator
kubectl apply -f deploy/crds/cache.example.com_v1alpha1_memcached_cr.yaml
```

```bash
# run as a deployment inside of the cluster
```

```bash
# run with olm
```

## Random Scratch

```bash
# create project
mkdir over-ez-operator

cd over-ez-operator

operator-sdk init --domain=mydomain.com --repo=github.com/mydomain/over-ez-operator
```

```bash
# tree
.
├── bin
│   └── manager
├── config
│   ├── certmanager
│   │   ├── certificate.yaml
│   │   ├── kustomization.yaml
│   │   └── kustomizeconfig.yaml
│   ├── default
│   │   ├── kustomization.yaml
│   │   ├── manager_auth_proxy_patch.yaml
│   │   ├── manager_webhook_patch.yaml
│   │   └── webhookcainjection_patch.yaml
│   ├── manager
│   │   ├── kustomization.yaml
│   │   └── manager.yaml
│   ├── prometheus
│   │   ├── kustomization.yaml
│   │   └── monitor.yaml
│   ├── rbac
│   │   ├── auth_proxy_client_clusterrole.yaml
│   │   ├── auth_proxy_role_binding.yaml
│   │   ├── auth_proxy_role.yaml
│   │   ├── auth_proxy_service.yaml
│   │   ├── kustomization.yaml
│   │   ├── leader_election_role_binding.yaml
│   │   ├── leader_election_role.yaml
│   │   └── role_binding.yaml
│   ├── scorecard
│   │   ├── bases
│   │   │   └── config.yaml
│   │   ├── kustomization.yaml
│   │   └── patches
│   │       ├── basic.config.yaml
│   │       └── olm.config.yaml
│   └── webhook
│       ├── kustomization.yaml
│       ├── kustomizeconfig.yaml
│       └── service.yaml
├── Dockerfile
├── go.mod
├── go.sum
├── hack
│   └── boilerplate.go.txt
├── main.go
├── Makefile
└── PROJECT
```

```bash
# create new api and controller
operator-sdk create api --group=lab00 --version=v1alpha1 --kind=Mycrd
Create Resource [y/n]
y
Create Controller [y/n]
y

make generate
make manifests
```

```bash
# create sample pod
kubectl create ns golang-lab00
```


