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

package main

import (
	"encoding/base64"
	"flag"
	"github.com/tektoncd/pipeline/pkg/logging"
	"io/ioutil"
	"os"
	"strings"
)

const (
	prFile = "prmetadata.json"
)

var (
	prName    = flag.String("prname", "", "The name of the pipelinerun which needs to be resolved")
	namespace = flag.String("namespace", "", "The namespace to fetch the resource from")
)

func main() {
	flag.Parse()
	logger, _ := logging.NewLogger("", "metadata-init")
	defer logger.Sync()

	prb64 := strings.TrimSpace(os.Getenv("PRMETADATA"))

	pr, err := base64.StdEncoding.DecodeString(prb64)

	if err != nil {
		logger.Fatalf("Could not decode pipelinerun metadata: %w", err)
	}

	err = ioutil.WriteFile(prFile, pr, 0644)

	if err != nil {
		logger.Fatalf("Error writing metadata file : %w", err)
	}
}
