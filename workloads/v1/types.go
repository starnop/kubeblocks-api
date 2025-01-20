/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

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
)

type ClusterComponentVolumeClaimTemplate struct {
	// Refers to the name of a volumeMount defined in either:
	//
	// - `componentDefinition.spec.runtime.containers[*].volumeMounts`
	// - `clusterDefinition.spec.componentDefs[*].podSpec.containers[*].volumeMounts` (deprecated)
	//
	// The value of `name` must match the `name` field of a volumeMount specified in the corresponding `volumeMounts` array.
	//
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Defines the desired characteristics of a PersistentVolumeClaim that will be created for the volume
	// with the mount name specified in the `name` field.
	//
	// When a Pod is created for this ClusterComponent, a new PVC will be created based on the specification
	// defined in the `spec` field. The PVC will be associated with the volume mount specified by the `name` field.
	//
	// +optional
	Spec PersistentVolumeClaimSpec `json:"spec,omitempty"`
}

type PersistentVolumeClaimSpec struct {
	// Contains the desired access modes the volume should have.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1.
	//
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	AccessModes []corev1.PersistentVolumeAccessMode `json:"accessModes,omitempty" protobuf:"bytes,1,rep,name=accessModes,casttype=PersistentVolumeAccessMode"`

	// Represents the minimum resources the volume should have.
	// If the RecoverVolumeExpansionFailure feature is enabled, users are allowed to specify resource requirements that
	// are lower than the previous value but must still be higher than the capacity recorded in the status field of the claim.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources.
	//
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	// TODO: 需要解决 go mod 冲突
	//Resources corev1.VolumeResourceRequirements `json:"resources,omitempty" protobuf:"bytes,2,opt,name=resources"`

	// The name of the StorageClass required by the claim.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1.
	//
	// +optional
	StorageClassName *string `json:"storageClassName,omitempty" protobuf:"bytes,5,opt,name=storageClassName"`

	// Defines what type of volume is required by the claim, either Block or Filesystem.
	//
	// +optional
	VolumeMode *corev1.PersistentVolumeMode `json:"volumeMode,omitempty" protobuf:"bytes,6,opt,name=volumeMode,casttype=PersistentVolumeMode"`

	// volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim.
	//
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#volumeattributesclass
	//
	// +optional
	VolumeAttributesClassName *string `json:"volumeAttributesClassName,omitempty"`
}

type SchedulingPolicy struct {
	// If specified, the Pod will be dispatched by specified scheduler.
	// If not specified, the Pod will be dispatched by default scheduler.
	//
	// +optional
	SchedulerName string `json:"schedulerName,omitempty"`

	// NodeSelector is a selector which must be true for the Pod to fit on a node.
	// Selector which must match a node's labels for the Pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	//
	// +optional
	// +mapType=atomic
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// NodeName is a request to schedule this Pod onto a specific node. If it is non-empty,
	// the scheduler simply schedules this Pod onto that node, assuming that it fits resource
	// requirements.
	//
	// +optional
	NodeName string `json:"nodeName,omitempty"`

	// Specifies a group of affinity scheduling rules of the Cluster, including NodeAffinity, PodAffinity, and PodAntiAffinity.
	//
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// Allows Pods to be scheduled onto nodes with matching taints.
	// Each toleration in the array allows the Pod to tolerate node taints based on
	// specified `key`, `value`, `effect`, and `operator`.
	//
	// - The `key`, `value`, and `effect` identify the taint that the toleration matches.
	// - The `operator` determines how the toleration matches the taint.
	//
	// Pods with matching tolerations are allowed to be scheduled on tainted nodes, typically reserved for specific purposes.
	//
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// TopologySpreadConstraints describes how a group of Pods ought to spread across topology
	// domains. Scheduler will schedule Pods in a way which abides by the constraints.
	// All topologySpreadConstraints are ANDed.
	//
	// +optional
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}

// InstanceTemplate allows customization of individual replica configurations in a Component.
type InstanceTemplate struct {
	// Name specifies the unique name of the instance Pod created using this InstanceTemplate.
	// This name is constructed by concatenating the Component's name, the template's name, and the instance's ordinal
	// using the pattern: $(cluster.name)-$(component.name)-$(template.name)-$(ordinal). Ordinals start from 0.
	// The specified name overrides any default naming conventions or patterns.
	//
	// +kubebuilder:validation:MaxLength=54
	// +kubebuilder:validation:Pattern:=`^[a-z0-9]([a-z0-9\.\-]*[a-z0-9])?$`
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// Specifies the number of instances (Pods) to create from this InstanceTemplate.
	// This field allows setting how many replicated instances of the Component,
	// with the specific overrides in the InstanceTemplate, are created.
	// The default value is 1. A value of 0 disables instance creation.
	//
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// Specifies the desired Ordinals of this InstanceTemplate.
	// The Ordinals used to specify the ordinal of the instance (pod) names to be generated under this InstanceTemplate.
	//
	// For example, if Ordinals is {ranges: [{start: 0, end: 1}], discrete: [7]},
	// then the instance names generated under this InstanceTemplate would be
	// $(cluster.name)-$(component.name)-$(template.name)-0、$(cluster.name)-$(component.name)-$(template.name)-1 and
	// $(cluster.name)-$(component.name)-$(template.name)-7
	Ordinals Ordinals `json:"ordinals,omitempty"`

	// Specifies a map of key-value pairs to be merged into the Pod's existing annotations.
	// Existing keys will have their values overwritten, while new keys will be added to the annotations.
	//
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Specifies a map of key-value pairs that will be merged into the Pod's existing labels.
	// Values for existing keys will be overwritten, and new keys will be added.
	//
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Specifies an override for the first container's image in the Pod.
	//
	// +optional
	Image *string `json:"image,omitempty"`

	// Specifies the scheduling policy for the Component.
	//
	// +optional
	SchedulingPolicy *SchedulingPolicy `json:"schedulingPolicy,omitempty"`

	// Specifies an override for the resource requirements of the first container in the Pod.
	// This field allows for customizing resource allocation (CPU, memory, etc.) for the container.
	//
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// Defines Env to override.
	// Add new or override existing envs.
	// +optional
	Env []corev1.EnvVar `json:"env,omitempty"`

	// Defines Volumes to override.
	// Add new or override existing volumes.
	// +optional
	Volumes []corev1.Volume `json:"volumes,omitempty"`

	// Defines VolumeMounts to override.
	// Add new or override existing volume mounts of the first container in the Pod.
	// +optional
	VolumeMounts []corev1.VolumeMount `json:"volumeMounts,omitempty"`

	// Defines VolumeClaimTemplates to override.
	// Add new or override existing volume claim templates.
	// +optional
	VolumeClaimTemplates []ClusterComponentVolumeClaimTemplate `json:"volumeClaimTemplates,omitempty"`
}

// Range represents a range with a start and an end value.
// It is used to define a continuous segment.
type Range struct {
	Start int32 `json:"start"`
	End   int32 `json:"end"`
}

// Ordinals represents a combination of continuous segments and individual values.
type Ordinals struct {
	Ranges   []Range `json:"ranges,omitempty"`
	Discrete []int32 `json:"discrete,omitempty"`
}
