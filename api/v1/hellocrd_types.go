/*
Copyright 2022.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HellocrdSpec defines the desired state of Hellocrd
type HellocrdSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Hellocrd. Edit hellocrd_types.go to remove/update
	Foo string `json:"foo,omitempty"`
	// hello crd spec
	ContainerImageNamespace string `json:"containerImageNamespace,omitempty"`
	ContainerImage          string `json:"containerImage,omitempty"`
	ContainerTag            string `json:"containerTag,omitempty"`
}

// HellocrdStatus defines the observed state of Hellocrd
type HellocrdStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	HelloStatus string `json:"helloStatus,omitempty"`
	LastPodName string `json:"lastPodName,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Hellocrd is the Schema for the hellocrds API
type Hellocrd struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HellocrdSpec   `json:"spec,omitempty"`
	Status HellocrdStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HellocrdList contains a list of Hellocrd
type HellocrdList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Hellocrd `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Hellocrd{}, &HellocrdList{})
}
