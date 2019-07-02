/*
 *
 * Copyright 2019 The Tekton Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * /
 */

package v1alpha1

import (
	"github.com/knative/pkg/apis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Add validation for TaskConditions?
type TaskCondition struct {
	ConditionRef string `json:"conditionRef"`
	// TODO: Support a ConditionSpec?
	// +optional
	Params []Param `json:"params,omitempty"`
}

// Check that Task may be validated and defaulted.
var _ apis.Validatable = (*Condition)(nil)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Task represents a collection of sequential steps that are run as part of a
// Pipeline using a set of inputs and producing a set of outputs. Tasks execute
// when TaskRuns are created that provide the input parameters and resources and
// output resources the Task requires.
//
// +k8s:openapi-gen=true
type Condition struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Spec holds the desired state of the Condition from the client
	// +optional
	Spec ConditionSpec `json:"spec"`
}

type ConditionSpec struct {
	// +optional
	Params []ParamSpec `json:"params,omitempty"`
	// Check is a container whose exit code determines where a condition is true or false
	Check corev1.Container `json:"check,omitempty"`
}

type ConditionCheck TaskRun

type ConditionCheckStatus TaskRunStatus

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConditionList contains a list of Conditions
type ConditionList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Condition `json:"items"`
}

func NewConditionCheck(tr *TaskRun) *ConditionCheck {
	if tr == nil {
		return nil
	}

	cc := ConditionCheck(*tr)
	return &cc
}