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
	"github.com/ocklin/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	"github.com/ocklin/ndb-operator/pkg/resources"
	"k8s.io/apimachinery/pkg/api/errors"
	coreinformers "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"

	corev1 "k8s.io/api/core/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
)

type ConfigMapControlInterface interface {
	EnsureConfigMap(ndb *v1alpha1.Ndb) (*corev1.ConfigMap, error)
	UpdateConfigMap(ndb *v1alpha1.Ndb) (*corev1.ConfigMap, error)
	DeleteConfigMap(ndb *v1alpha1.Ndb) error
}

type ConfigMapControl struct {
	ConfigMapControlInterface

	k8client kubernetes.Interface

	configMapLister       corelisters.ConfigMapLister
	configMapListerSynced cache.InformerSynced
}

// NewConfigMapControl creates a new ConfigMapControl
func NewConfigMapControl(client kubernetes.Interface,
	configMapInformer coreinformers.ConfigMapInformer) ConfigMapControlInterface {

	configMapControl := &ConfigMapControl{
		k8client:              client,
		configMapLister:       configMapInformer.Lister(),
		configMapListerSynced: configMapInformer.Informer().HasSynced,
	}

	return configMapControl
}

func (rcmc *ConfigMapControl) EnsureConfigMap(ndb *v1alpha1.Ndb) (*corev1.ConfigMap, error) {

	// Get the StatefulSet with the name specified in Ndb.spec
	// TODO: probably dangerous if we do not sync caches, maybe fetch uncached directly
	cm, err := rcmc.configMapLister.ConfigMaps(ndb.Namespace).Get(ndb.GetConfigMapName())

	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		klog.Infof("Creating ConfigMap %s/%s", ndb.Namespace, ndb.GetConfigMapName())

		cm = resources.GenerateConfigMapObject(ndb)
		cm, err = rcmc.k8client.CoreV1().ConfigMaps(ndb.Namespace).Create(cm)
	}
	return cm, err
}

func (rcmc *ConfigMapControl) UpdateConfigMap(ndb *v1alpha1.Ndb) (*corev1.ConfigMap, error) {
	return nil, nil
}

func (rcmc *ConfigMapControl) DeleteConfigMap(ndb *v1alpha1.Ndb) error {
	return nil
}
