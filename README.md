# ndb-operator

This repository implements a simple controller for watching Ndb resources as
defined with a CustomResourceDefinition (CRD). 

This is pre-alpha - mostly prototyping for testing concepts.

## Details

The ndb controller uses [client-go library](https://github.com/kubernetes/client-go/tree/master/tools/cache) extensively.

## Fetch ndb-operator and its dependencies

The ndb-operator uses go modules and has been developed with go 1.13. 

```sh
git clone <this repo>
cd ndb-operator
```

## Changes to the types 

Note, if you intend to change the types then you will need to 
[generate code](#changes-to-the-types) which only seems to 
support old style $GOPATH. 

First it depends on the vendor directory. Populate that with

```sh
go mod vendor
```

The [code generator script](hack/update-codegen.sh) script will generate files &
directories:

* `pkg/apis/ndbcontroller/v1alpha1/zz_generated.deepcopy.go`
* `pkg/generated/`

It requires some copying as described in the script.

Generators are in [k8s.io/code-generator](https://github.com/kubernetes/code-generator)
and generate a typed client, informers, listers and deep-copy functions.

## Building

You may have to set the `OS` variable to your operating system in the Makefile.

```sh
# Build ndb-operator 
make build
```

## Docker image building

You can build your own ndb cluster images but you don't have to. Currently public image 8.0.22 is used.

**Prerequisite**: You have a build directory available that is build for OL8. 
You can use a OL8 build-container in docker/ol8 for that or download a readily compiled OL8 build.

If you use minikube then set the environment to minikube first before building the image.

```sh
# copy necessary ndb binaries 
make install-minimal

# point to minikube
$ eval $(minikube docker-env)

# build docker image
make build-docker
```

## Running

**Prerequisite**: operator built, docker images built and made available in kubernetes 

Create custom resource definitions and the roles:

```sh
NAMESPACE=default
kubectl -n ${NAMESPACE} apply -f artifacts/manifests/crd.yaml
sed -e "s/<NAMESPACE>/${NAMESPACE}/g" artifacts/manifests/rbac.yaml | kubectl -n ${NAMESPACE} apply -f -
```

or use helm

```sh
helm install ndb-operator helm
```

the proceed with starting the operator - currently outside of the k8 cluster:

```sh
./ndb-operator -kubeconfig=$HOME/.kube/config

# create a custom resource of type Ndb
kubectl apply -f artifacts/examples/example-ndb.yaml

# check statefulsets created through the custom resource
kubectl get pods,statefulsets

# watch pods change state
kubectl get pods -w

# "log into" pods with 
kubectl exec -ti pod/example-ndb-mgmd-0 -- /bin/bash
```

You can delete the cluster installation again with


```sh
kubectl delete -f artifacts/examples/example-ndb.yaml
```


## Defining types

The Ndb instance of our custom resource has an attached Spec, 
which is defined in [types.go](pkg/apis/ndbcontroller/types.go)

## Validation

To validate custom resources, use the [`CustomResourceValidation`](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/#validation) feature.

This feature is beta and enabled by default in v1.9.

### Example

The schema in [`crd-validation.yaml`](./artifacts/examples/crd-validation.yaml) applies the following validation on the custom resource:
`spec.replicas` must be an integer and must have a minimum value of 1 and a maximum value of 10.

In the above steps, use `crd-validation.yaml` to create the CRD:

```sh
# create a CustomResourceDefinition supporting validation
kubectl create -f artifacts/examples/crd-validation.yaml
```

## Subresources

TBD

## Cleanup

You can clean up the created CustomResourceDefinition with:

    kubectl delete crd ndbs.ndbcontroller.k8s.io


