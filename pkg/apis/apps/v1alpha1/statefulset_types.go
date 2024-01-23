/*
Copyright 2023 The Karmada Authors.

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

package v1alpha1

import (
	policy "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// ResourceKindCrossClusterService is kind name of CrossClusterService.
	ResourceKindCrossClusterStatefulset            = "CrossClusterStatefulset"
	ResourceSingularCrossClusterStatefulset        = "CrossClusterstatefulset"
	ResourcePluralCrossClusterStatefulset          = "CrossClusterStatefulsets"
	ResourceNamespaceScopedCrossClusterStatefulset = true
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=ksts,categories={karmada-io}

type CrossClusterStatefulSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the CrossClusterStatefulSet.
	Spec CrossClusterStatefulSetSpec `json:"spec,omitempty"`
}

// CrossClusterStatefulSetSpec is the desired state of the CrossClusterService.
type CrossClusterStatefulSetSpec struct {
	ResourceSelector policy.ResourceSelector `json:"resourceSelector,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CrossClusterStatefulSetList contains a list of CrossClusterStatefulSet.
type CrossClusterStatefulSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// Items           []FederatedResourceQuota `json:"items"`
}
