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

package v1alpha1

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PreflightCheckSpec defines the desired state of PreflightCheck
type PreflightCheckSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PreflightCheck. Edit preflightcheck_types.go to remove/update
	Foo string `json:"foo,omitempty"`

	// Image is the container image to run preflight against.
	Image string `json:"image"`

	// LogLevel represents the preflight log level.
	// +kubebuilder:validation:Enum=Info;Warn;Debug;Trace
	LogLevel *string `json:"logLevel,omitempty"`

	// DockerConfigSecretRef is a secret containing a key config.json with a dockerconfig.json
	// as its contents.
	DockerConfigSecretRef *string `json:"dockerConfigSecretRef,omitempty"`

	// PreflightImage overrides the default preflight stable container image.
	PreflightImage *string `json:"preflightImage,omitempty"`

	CheckOptions CheckOptions `json:"checkOptions"`
}

// +kubebuilder:validation:MaxProperties:=1
type CheckOptions struct {
	ContainerOptions *ContainerOptions `json:"containerOptions,omitempty"`
	OperatorOptions  *OperatorOptions  `json:"operatorOptions,omitempty"`
}

type ContainerOptions struct {
	CertificationProjectID *string `json:"certificationProjectID,omitempty"`
	PyxisAPITokenSecretRef *string `json:"pyxisAPITokenSecretRef,omitempty"`
}

type OperatorOptions struct {
	KubeconfigSecretRef     string `json:"kubeconfigSecretRef"`
	IndexImage              string `json:"indexImage"`
	ScorecardServiceAccount string `json:"scorecardServiceAccount,omitempty"`
	ScorecardNamespace      string `json:"scorecardNamespace,omitempty"`
	ScorecardWaitTime       string `json:"scorecardWaitTime,omitempty"`
	DeploymentChannel       string `json:"deploymentChannel,omitempty"`
}

func (pc *PreflightCheck) GenerateJob() batchv1.Job {
	job := PreflightCheckJobGenerator{*pc}
	return job.Generate()
}

// PreflightCheckStatus defines the observed state of PreflightCheck
type PreflightCheckStatus struct {
	Jobs      []string `json:"jobs,omitempty"`
	Completed bool     `json:"completed"`
	Type      string   `json:"type,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.image"
//+kubebuilder:printcolumn:name="Type",type="string",JSONPath=".status.type"
//+kubebuilder:printcolumn:name="Successful",type="boolean",JSONPath=".status.completed"
//+kubebuilder:printcolumn:name="Job",type="string",JSONPath=".status.jobs[0]"

// PreflightCheck is the Schema for the preflightchecks API
type PreflightCheck struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PreflightCheckSpec   `json:"spec,omitempty"`
	Status PreflightCheckStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PreflightCheckList contains a list of PreflightCheck
type PreflightCheckList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PreflightCheck `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PreflightCheck{}, &PreflightCheckList{})
}

func (pc *PreflightCheck) CheckType() string {
	t := "container"

	if pc.Spec.CheckOptions.OperatorOptions != nil {
		t = "operator"
	}

	return t
}
