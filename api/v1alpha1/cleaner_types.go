/*
Copyright 2023. projectsveltos.io. All rights reserved.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	libsveltosv1alpha1 "github.com/projectsveltos/libsveltos/api/v1alpha1"
)

// Action specifies the action to take on matching resources
// +kubebuilder:validation:Enum:=Delete;Transform
type Action string

const (
	// ActionDelete will delete the resource
	ActionDelete = Action("Delete")

	// ActionTransform will update object
	ActionTransform = Action("Transform")
)

const (
	// CleanerFinalizer allows Reconciler to clean up resources associated with
	// Cleaner instance before removing it from the apiserver.
	CleanerFinalizer = "cleanerfinalizer.projectsveltos.io"
)

type ResourceSelector struct {
	// Namespace of the resource deployed in the  Cluster.
	// Empty for resources scoped at cluster level.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// Group of the resource deployed in the Cluster.
	Group string `json:"group"`

	// Version of the resource deployed in the Cluster.
	Version string `json:"version"`

	// Kind of the resource deployed in the Cluster.
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind"`

	// LabelFilters allows to filter resources based on current labels.
	LabelFilters []libsveltosv1alpha1.LabelFilter `json:"labelFilters,omitempty"`

	// Evaluate contains a function "evaluate" in lua language.
	// The function will be passed one of the object selected based on
	// above criteria.
	// Must return struct with field "matching"
	// representing whether object is a match.
	// +optional
	Evaluate string `json:"evaluate,omitempty"`
}

type ResourcePolicySet struct {
	// ResourceSelectors identifies what resources to select
	ResourceSelectors []ResourceSelector `json:"resourceSelectors"`

	// This field is optional and can be used to specify a Lua function
	// that will be used to further select a subset of the resources that
	// have already been selected using the ResourceSelector field.
	// The function will receive the array of resources selected by ResourceSelectors.
	// If this field is not specified, all resources selected by the ResourceSelector
	// field will be considered.
	// This field allows to perform more complex filtering or selection operations
	// on the resources, looking at all resources together.
	// This can be useful for more sophisticated tasks, such as identifying resources
	// that are related to each other or that have similar properties.
	AggregatedSelection string `json:"aggregatedSelection,omitempty"`
}

// NotificationType specifies different type of notifications
// +kubebuilder:validation:Enum:=CleanerReport;Slack;Webex;Discord
type NotificationType string

const (
	// NotificationTypeCleanerReport refers to generating a CleanerReport instance
	NotificationTypeCleanerReport = NotificationType("CleanerReport")

	// NotificationTypeSlack refers to generating a Slack message
	NotificationTypeSlack = NotificationType("Slack")

	// NotificationTypeWebex refers to generating a Webex message
	NotificationTypeWebex = NotificationType("Webex")

	// NotificationTypeDiscord refers to generating a Discord message
	NotificationTypeDiscord = NotificationType("Discord")
)

type Notification struct {
	// Name of the notification check.
	// Must be a DNS_LABEL and unique within the Cleaner.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// NotificationType specifies the type of notification
	Type NotificationType `json:"type"`

	// NotificationRef is a reference to a notification-specific resource that holds
	// the details for the notification.
	// +optional
	NotificationRef *corev1.ObjectReference `json:"notificationRef,omitempty"`
}

// CleanerSpec defines the desired state of Cleaner
type CleanerSpec struct {
	// ResourcePolicySet identifies a group of resources
	ResourcePolicySet ResourcePolicySet `json:"resourcePolicySet"`

	// Action indicates the action to take on selected object. Default action
	// is to delete object. If set to transform, the transform function
	// will be invoked and then object will be updated.
	// +kubebuilder:default:=Delete
	Action Action `json:"action,omitempty"`

	// Transform contains a function "transform" in lua language.
	// When Action is set to *Transform*, this function will be invoked
	// and be passed one of the object selected based on
	// above criteria.
	// Must the new object that will be applied
	// +optional
	Transform string `json:"transform,omitempty"`

	// DryRun if set to true, will have controller delete/update no resource.
	// All matching resources will be listed in logs
	// +kubebuilder:default:=false
	// +optional
	DryRun bool `json:"dryRun,omitempty"`

	// Schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	Schedule string `json:"schedule"`

	// Optional deadline in seconds for starting the job if it misses scheduled
	// time for any reason.  Missed jobs executions will be counted as failed ones.
	// +optional
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`

	// Notification is a list of source of events to evaluate.
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	// +optional
	Notifications []Notification `json:"notifications,omitempty"`
}

// CleanerStatus defines the observed state of Cleaner
type CleanerStatus struct {
	// Information when next snapshot is scheduled
	// +optional
	NextScheduleTime *metav1.Time `json:"nextScheduleTime,omitempty"`

	// Information when was the last time a snapshot was successfully scheduled.
	// +optional
	LastRunTime *metav1.Time `json:"lastRunTime,omitempty"`

	// FailureMessage provides more information about the error, if
	// any occurred
	FailureMessage *string `json:"failureMessage,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:path=cleaners,scope=Cluster
//+kubebuilder:subresource:status

// Cleaner is the Schema for the cleaners API
type Cleaner struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CleanerSpec   `json:"spec,omitempty"`
	Status CleanerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CleanerList contains a list of Cleaner
type CleanerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cleaner `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cleaner{}, &CleanerList{})
}
