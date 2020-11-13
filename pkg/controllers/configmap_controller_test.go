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

// NOTE: NOT READY AT ALL - FIX BUT DON'T USE

package controllers

import (
	"encoding/json"
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
)

func Test_TestThingsrelatedToConfigMaps(t *testing.T) {

	configString := "Version 1"
	ns := metav1.NamespaceDefault

	d1 := map[string]string{
		"config.ini": configString,
	}
	d2 := map[string]string{
		"config.ini": "Version 2",
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "configtest",
			Namespace: ns,
		},
		Data: d1,
	}

	cm2 := cm.DeepCopy()

	j, err := json.Marshal(cm)
	if err != nil {
		t.Error(err)
	}

	cm2.Data = d2

	j2, err := json.Marshal(cm2)
	if err != nil {
		t.Error(err)
	}
	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(
		j, j2, corev1.ConfigMap{})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(patchBytes))

	t.Fail()
}

func TestCreateConfigMap(t *testing.T) {

	f := newFixture(t)
	defer f.close()

	ns := metav1.NamespaceDefault
	ndb := newNdb(ns, "test", 1)
	ndb.Spec.Mysqld.NodeCount = int32Ptr(7)

	// we first need to set up arrays with objects ...
	f.ndbLister = append(f.ndbLister, ndb)
	f.objects = append(f.objects, ndb)

	// ... before we init the fake clients with those objects.
	// objects not listed in arrays at fakeclient setup will eventually be deleted
	f.init()

	cmc := NewConfigMapControl(f.kubeclient, f.k8If.Core().V1().ConfigMaps())

	f.start()

	cm, err := cmc.EnsureConfigMap(ndb)

	if err != nil {
		t.Errorf("Unexpected error EnsuringConfigMap: %v", err)
	}
	if cm == nil {
		t.Errorf("Unexpected error EnsuringConfigMap: return null pointer")
	}

	f.expectCreateAction(ndb.GetNamespace(), "configmap", cm)

	rcmc := cmc.(*ConfigMapControl)

	// Wait for the caches to be synced before using Lister to get new config map
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(f.stopCh, rcmc.configMapListerSynced); !ok {
		t.Errorf("failed to wait for caches to sync")
		return
	}

	// Get the StatefulSet with the name specified in Ndb.spec
	cmget, err := rcmc.configMapLister.ConfigMaps(ndb.Namespace).Get(ndb.GetConfigMapName())
	if err != nil {
		t.Errorf("Unexpected error getting created ConfigMap: %v", err)
	}
	if cmget == nil {
		t.Errorf("Unexpected error EnsuringConfigMap: didn't find created ConfigMap")
	}

	ndb.Spec.Mysqld.NodeCount = int32Ptr(12)

	cm, err = cmc.PatchConfigMap(ndb)

	s, _ := json.MarshalIndent(cm, "", "  ")

	t.Errorf(string(s))
}
