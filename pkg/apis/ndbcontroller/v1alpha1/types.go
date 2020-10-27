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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// Ndb is a specification for a Ndb resource
type Ndb struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NdbSpec   `json:"spec"`
	Status NdbStatus `json:"status"`
}

type NdbNdbdSpec struct {
	NoOfReplicas *int32 `json:"noofreplicas"`
	NodeCount    *int32 `json:"nodecount"`
	Name         string `json:"deploymentName"`
}

type NdbMgmdSpec struct {
	NodeCount *int32 `json:"nodecount"`
	Name      string `json:"name"`
}

type NdbMysqldSpec struct {
	NodeCount *int32 `json:"nodecount"`
	Name      string `json:"name"`
}

// NdbSpec is the spec for a Ndb resource
type NdbSpec struct {
	DeploymentName string        `json:"deploymentName"`
	Mgmd           NdbMgmdSpec   `json:"mgmd"`
	Ndbd           NdbNdbdSpec   `json:"ndbd"`
	Mysqld         NdbMysqldSpec `json:"mysqld"`

	// Config allows a user to specify a custom configuration file for MySQL.
	// +optional
	Config *corev1.LocalObjectReference `json:"config,omitempty"`
}

// NdbStatus is the status for a Ndb resource
type NdbStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NdbList is a list of Ndb resources
type NdbList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Ndb `json:"items"`
}
