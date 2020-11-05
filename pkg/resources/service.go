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
func NewService(ndb *v1alpha1.Ndb, mgmd bool, svcName string) *corev1.Service {

	selector := ndb.GetDataNodeLabels()
	svcType := corev1.ServiceTypeClusterIP
	clusterIP := corev1.ClusterIPNone

	if mgmd {
		selector = ndb.GetManagementNodeLabels()
		svcType = corev1.ServiceTypeLoadBalancer
		clusterIP = ""
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Labels: ndb.GetLabels(),
			Name:   svcName,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(ndb, schema.GroupVersionKind{
					Group:   v1.SchemeGroupVersion.Group,
					Version: v1.SchemeGroupVersion.Version,
					Kind:    "Ndb",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			PublishNotReadyAddresses: true,
			Ports: []corev1.ServicePort{
				// TODO: two ports in array didn't work, at least not exposing via minikube tunnel
				corev1.ServicePort{Port: 8080, Name: "agent", Protocol: "TCP"},
			},
			Selector:  selector,
			ClusterIP: clusterIP,
			Type:      svcType,
		},
	}

	return svc
}
