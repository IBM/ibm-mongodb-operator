//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Image struct {
	Tag string `json:"tag,omitempty"`
}

// MongoDBSpec defines the desired state of MongoDB
type MongoDBSpec struct {
	ImageRegistry  string `json:"imageRegistry,omitempty"`
	Replicas       int    `json:"replicas,omitempty"`
	StorageClass   string `json:"storageClass,omitempty"`
	InitImage      Image  `json:"initImage,omitempty"`
	BootstrapImage Image  `json:"bootstrapImage,omitempty"`
	MetricsImage   Image  `json:"metricsImage,omitempty"`
	Resources      corev1.ResourceRequirements `json:"resources,omitempty"`
	// CpuLimit			 	string `json:"resources.limits.cpu,omitempty"`
	// MemoryLimit			string `json:"resources.limits.memory,omitempty"`
	// CpuRequest			string `json:"resources.requests.cpu,omitempty"`
	// MemoryRequest		string `json:"resources.requests.memory,omitempty"`
}

// MongoDBStatus defines the observed state of MongoDB
type MongoDBStatus struct {
	StorageClass string `json:"storageClass,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDB is the Schema for the mongodbs API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=mongodbs,scope=Namespaced
type MongoDB struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MongoDBSpec   `json:"spec,omitempty"`
	Status MongoDBStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MongoDBList contains a list of MongoDB
type MongoDBList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MongoDB `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MongoDB{}, &MongoDBList{})
}
