/*
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
 *
 */

package resources

import (
	"encoding/base64"
	"encoding/json"
	"flag"

	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/tektoncd/pipeline/pkg/names"
	"golang.org/x/xerrors"
	corev1 "k8s.io/api/core/v1"
)

var (
	metadataPrefix = "metadata-init"
	// The container with Git that we use to implement the Git source step.
	metadataImage = flag.String("metadata-image", "metadata-init:latest",
		"The container image containing our Metadata binary.")
)

// ConditionMetadata represents the metadata about pipelinerun state that
// is accessible to a condition
type ConditionMetadata struct {
	PipelineRun *v1alpha1.PipelineRun `json:"pipelinerun,omitempty"`
	Pipeline    *v1alpha1.Pipeline    `Json:"pipeline,omitempty`
	Resources   []*v1alpha1.PipelineResource
}

func getMetadataContainerSpec(run *v1alpha1.PipelineRun, rprt *ResolvedPipelineRunTask) (corev1.Container, error) {
	m := &ConditionMetadata{
		PipelineRun: run,
	}

	prb64, err := pipelineRunAsBase64(m)
	if err != nil {
		return corev1.Container{}, err
	}
	step := corev1.Container{
		Name:    names.SimpleNameGenerator.RestrictLengthWithRandomSuffix(metadataPrefix),
		Image:   *metadataImage,
		Command: []string{"/ko-app/metadata-init"},
		Env: []corev1.EnvVar{{
			Name:  "PRMETADATA",
			Value: prb64,
		}},
		WorkingDir: "/workspace/metadata/",
	}

	return step, nil
}

func pipelineRunAsBase64(m *ConditionMetadata) (string, error) {
	asJSON, err := json.Marshal(m)
	if err != nil {
		return "", xerrors.Errorf("Could not marshall PipelineRun %s to JSON: %w", m.PipelineRun.Name, err)
	}

	return base64.StdEncoding.EncodeToString(asJSON), nil
}
