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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TargetDeployment struct {
	Name      string `json:"name"`
	APIPrefix string `json:"apiPrefix"`
	Port      string `json:"port"`
	Service   string `json:"service"`
}

// CustomAutoscalerSpec определяет желаемое состояние CustomAutoscaler
type CustomAutoscalerSpec struct {
	TargetDeployments     []TargetDeployment `json:"targetDeployments"`
	IdleTime              int64              `json:"idleTime"`
	CpuThreshold          int64              `json:"cpuThreshold"`
	MemoryThreshold       int64              `json:"memoryThreshold"`
	ResponseTimeThreshold int64              `json:"responseTimeThreshold"`
}

// CustomAutoscalerStatus defines the observed state of CustomAutoscaler
type CustomAutoscalerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CustomAutoscaler is the Schema for the customautoscalers API
type CustomAutoscaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomAutoscalerSpec   `json:"spec,omitempty"`
	Status CustomAutoscalerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CustomAutoscalerList contains a list of CustomAutoscaler
type CustomAutoscalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomAutoscaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CustomAutoscaler{}, &CustomAutoscalerList{})
}
