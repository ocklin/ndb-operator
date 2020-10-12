// Copyright 2020 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"encoding/json"

	"github.com/ocklin/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	"github.com/ocklin/ndb-operator/pkg/resources"
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/klog"
)

/* StatefulSetControlInterface defines the interface that the
   wraps around the creation and update of StatefulSets for node types. It
   is implemented as an interface to enable testing. */

type StatefulSetControlInterface interface {
	EnsureStatefulSet(ndb *v1alpha1.Ndb) (*apps.StatefulSet, error)
	GetStatefulSet(ndb *v1alpha1.Ndb) (*apps.StatefulSet, error)
	Patch(ndb *v1alpha1.Ndb, old *apps.StatefulSet) (*apps.StatefulSet, error)
}

type realStatefulSetControl struct {
	client            kubernetes.Interface
	statefulSetLister appslisters.StatefulSetLister
	statefulSetType   resources.StatefulSetInterface
}

// NewRealStatefulSetControl creates a concrete implementation of the
// StatefulSetControlInterface.
func NewRealStatefulSetControl(client kubernetes.Interface, statefulSetLister appslisters.StatefulSetLister) StatefulSetControlInterface {
	return &realStatefulSetControl{client: client, statefulSetLister: statefulSetLister}
}

func (rssc *realStatefulSetControl) GetStatefulSet(ndb *v1alpha1.Ndb) (*apps.StatefulSet, error) {
	return rssc.statefulSetLister.StatefulSets(ndb.Namespace).Get(ndb.Name)
}

// PatchStatefulSet performs a direct patch update for the specified StatefulSet.
func patchStatefulSet(client kubernetes.Interface, oldData *apps.StatefulSet, newData *apps.StatefulSet) (*apps.StatefulSet, error) {
	originalJSON, err := json.Marshal(oldData)
	if err != nil {
		return nil, err
	}

	//klog.Infof("Patching StatefulSet old: %s", string(originalJSON))

	updatedJSON, err := json.Marshal(newData)
	if err != nil {
		return nil, err
	}

	//klog.Infof("Patching StatefulSet new: %s", string(updatedJSON))

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(
		originalJSON, updatedJSON, apps.StatefulSet{})
	if err != nil {
		return nil, err
	}
	klog.Infof("Patching StatefulSet %q: %s", types.NamespacedName{
		Namespace: oldData.Namespace,
		Name:      oldData.Name}, string(patchBytes))

	result, err := client.AppsV1().StatefulSets(oldData.Namespace).Patch(oldData.Name,
		types.StrategicMergePatchType,
		patchBytes)

	if err != nil {
		klog.Errorf("Failed to patch StatefulSet: %v", err)
		return nil, err
	}

	return result, nil
}

func (rssc *realStatefulSetControl) Patch(ndb *v1alpha1.Ndb, old *apps.StatefulSet) (*apps.StatefulSet, error) {

	oldCopy := old.DeepCopy()

	klog.Infof("Patch stateful set %s/%s Replicas: %d, DataNodes: %d",
		ndb.Namespace,
		rssc.statefulSetType.GetName(),
		*ndb.Spec.Ndbd.NoOfReplicas,
		*ndb.Spec.Ndbd.NodeCount)

	sfset := rssc.statefulSetType.NewStatefulSet(ndb)

	return patchStatefulSet(rssc.client, oldCopy, sfset)
}

func (rssc *realStatefulSetControl) EnsureStatefulSet(ndb *v1alpha1.Ndb) (*apps.StatefulSet, error) {

	// Get the StatefulSet with the name specified in Ndb.spec
	sfset, err := rssc.statefulSetLister.StatefulSets(ndb.Namespace).Get(rssc.statefulSetType.GetName())

	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		klog.Infof("Creating stateful set %s/%s Replicas: %d, DataNodes: %d",
			ndb.Namespace,
			rssc.statefulSetType.GetName(),
			*ndb.Spec.Ndbd.NoOfReplicas,
			*ndb.Spec.Ndbd.NodeCount)
		sfset = rssc.statefulSetType.NewStatefulSet(ndb)
		_, err = rssc.client.AppsV1().StatefulSets(ndb.Namespace).Create(sfset)
	}

	if err != nil {
		// re-queue if something went wrong
		klog.Errorf("Failed to create stateful set %s/%s replicas: %d with error: %s",
			ndb.Namespace, rssc.statefulSetType.GetName(),
			*ndb.Spec.Ndbd.NodeCount, err)
	}

	return sfset, err
}
