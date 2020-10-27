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

package resources

import (
	"github.com/ocklin/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
	"github.com/ocklin/ndb-operator/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const ConfigMapName = "ndb-config-ini"

func GenerateConfigMapObject(ndb *v1alpha1.Ndb) *corev1.ConfigMap {

	/*
		kind: ConfigMap
		apiVersion: v1
		metadata:
			name: config-ini
			namespace: default
			#uid: a7ec90d9-f2a9-11ea-95f5-000d3a2ebd7f
			#resourceVersion: '58127766'
			#creationTimestamp: '2020-09-09T14:35:06Z'
		data:
			config.ini: |
				[DB DEFAULT]
				....
	*/

	configStr, err := GetConfigString(ndb)
	if err != nil {
		return nil
	}

	data := map[string]string{
		"config.ini": configStr,
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ConfigMapName,
			Namespace: ndb.Namespace,
			Labels: map[string]string{
				constants.ClusterLabel: ndb.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(ndb, schema.GroupVersionKind{
					Group:   corev1.SchemeGroupVersion.Group,
					Version: corev1.SchemeGroupVersion.Version,
					Kind:    "Ndb",
				}),
			},
		},
		Data: data,
	}

	return cm
}
