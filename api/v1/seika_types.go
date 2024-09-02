/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SeikaSpec defines the desired state of Seika
type SeikaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Repurika map[string]int `json:"repurika,omitempty"`

	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	Template corev1.PodTemplateSpec `json:"template,omitempty"`
}

// SeikaStatus defines the observed state of Seika
type SeikaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Repurika map[string]int `json:"repurika,omitempty"`

	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`

	Nodes []string `json:"nodes,omitempty"`

	Ready string `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Nodes",type=string,JSONPath=`.status.nodes`
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.ready`
// Seika is the Schema for the seikas API
type Seika struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SeikaSpec   `json:"spec,omitempty"`
	Status SeikaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SeikaList contains a list of Seika
type SeikaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Seika `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Seika{}, &SeikaList{})
}
