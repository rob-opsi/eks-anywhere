// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetKubeVipAssets returns the eks a artifacts for kube-vip
func (r *ReleaseConfig) GetKubeVipAssets() ([]Artifact, error) {
	// Get Git tag for the project
	projectSource := "projects/plunder-app/kube-vip"
	tagFile := filepath.Join(r.BuildRepoSource, projectSource, "GIT_TAG")
	gitTag, err := readFile(tagFile)
	if err != nil {
		return nil, errors.Cause(err)
	}

	name := "kube-vip"
	repoName := fmt.Sprintf("plunder-app/%s", name)
	tagOptions := map[string]string{
		"gitTag": gitTag,
	}

	imageArtifact := &ImageArtifact{
		AssetName:       name,
		SourceImageURI:  r.GetSourceImageURI(name, repoName, tagOptions),
		ReleaseImageURI: r.GetReleaseImageURI(name, repoName, tagOptions),
		Arch:            []string{"amd64"},
		OS:              "linux",
		GitTag:          gitTag,
	}

	artifact := &Artifact{Image: imageArtifact}

	return []Artifact{*artifact}, nil
}
