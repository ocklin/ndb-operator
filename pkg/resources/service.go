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
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/ocklin/ndb-operator/pkg/apis/ndbcontroller/v1alpha1"
)

// NewForCluster will return a new headless Kubernetes service for a MySQL cluster
func NewService(ndb *v1alpha1.Ndb) *corev1.Service {
	mysqlPort := corev1.ServicePort{Port: 1186}
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: ndb.GetLabels(),
			Name:   ndb.GetServiceName(),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(ndb, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "Ndb",
				}),
			},
			Annotations: map[string]string{
				"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true",
			},
		},
		Spec: corev1.ServiceSpec{
			Ports:     []corev1.ServicePort{mysqlPort},
			Selector:  ndb.GetLabels(),
			ClusterIP: corev1.ClusterIPNone,
			Type:      corev1.ServiceTypeClusterIP,
			//Type: corev1.ServiceTypeNodePort,
		},
	}

	return svc
}
