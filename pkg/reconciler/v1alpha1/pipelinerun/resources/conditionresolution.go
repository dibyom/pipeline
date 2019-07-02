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
 */

package resources

import (
	"github.com/knative/pkg/apis"
	corev1 "k8s.io/api/core/v1"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
)

const (
	unnamedCheckNamePrefix = "condition-check-"
)
// GetCondition is a function used to retrieve PipelineConditions.
type GetCondition func(string) (*v1alpha1.Condition, error)

type ResolvedConditionCheck struct {
	ConditionCheckName string
	Condition          *v1alpha1.Condition
	ConditionCheck     *v1alpha1.ConditionCheck
}

type TaskConditionCheckState []*ResolvedConditionCheck

func (state TaskConditionCheckState) HasStarted() bool {
	hasStarted := true
	for _, j := range state {
		if j.ConditionCheck == nil {
			hasStarted = false
		}
	}
	return hasStarted
}

func (state TaskConditionCheckState) IsComplete() bool {
	if !state.HasStarted() {
		return false
	}
	isDone := true
	for _, rcc := range state {
		isDone = isDone && !rcc.ConditionCheck.Status.GetCondition(apis.ConditionSucceeded).IsUnknown()
	}
	return isDone
}

func (state TaskConditionCheckState) IsSuccess() bool {
	if !state.IsComplete() {
		return false
	}
	isSuccess := true
	for _, rcc := range state {
		isSuccess = isSuccess && rcc.ConditionCheck.Status.GetCondition(apis.ConditionSucceeded).IsTrue()
	}
	return isSuccess
}

// Convert a Condition to a TaskSpec
func (rcc *ResolvedConditionCheck) ConditionToTaskSpec() *v1alpha1.TaskSpec {
	// TODO(dibyom): Should be in SetDefaults?
	if rcc.Condition.Spec.Check.Name == "" {
		rcc.Condition.Spec.Check.Name = unnamedCheckNamePrefix + rcc.Condition.Name
	}

	t := &v1alpha1.TaskSpec{
		Steps: []corev1.Container{rcc.Condition.Spec.Check},
	}

	if len(rcc.Condition.Spec.Params) > 0 {
		t.Inputs = &v1alpha1.Inputs{
			Params: rcc.Condition.Spec.Params,
		}
	}

	return t
}

func (rcc *ResolvedConditionCheck) NewConditionCheckStatus() v1alpha1.ConditionCheckStatus {
	var checkStep corev1.ContainerState
	trs := rcc.ConditionCheck.Status
	for _, s := range trs.Steps {
		if s.Name == rcc.Condition.Spec.Check.Name {
			checkStep = s.ContainerState
			break
		}
	}

	return v1alpha1.ConditionCheckStatus{
		Status:         trs.Status,
		PodName:        trs.PodName,
		StartTime:      trs.StartTime,
		CompletionTime: trs.CompletionTime,
		Check:          checkStep,
	}
}