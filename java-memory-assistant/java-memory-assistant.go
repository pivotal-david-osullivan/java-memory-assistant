/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package java_memory_assistant

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type javaMemoryAssistant struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func JavaMemoryAssistant(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) (javaMemoryAssistant, libcnb.BOMEntry) {

	// Call libpak method to create a new 'contributor' which contributes our dependency to a 'Launch' layer
	contributor, entry := libpak.NewDependencyLayer(dependency, cache, libcnb.LayerTypes{
		Launch: true,
	})
	return javaMemoryAssistant{LayerContributor: contributor}, entry
}

func (w javaMemoryAssistant) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	w.LayerContributor.Logger = w.Logger

	return w.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {

		/*paths, err := filepath.Glob(artifact."java-memory-assistant-*.jar")
		if err != nil || len(paths) != 1 {
			return libcnb.Layer{}, fmt.Errorf("unable to locate agent jar\n%w", err)
		}*/
		agentPath := artifact.Name()

		// Create a bin directory so that the dependency is automatically added to $PATH at launch
		binDir := filepath.Join(layer.Path, "bin")

		if err := os.MkdirAll(binDir, 0755); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to mkdir\n%w", err)
		}

		if err := os.Symlink(filepath.Join(layer.Path, agentPath), filepath.Join(binDir, agentPath)); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to symlink agent\n%w", err)
		}
		// Finally add the agent to the JAVA_TOOL_OPTIONS env var via '-javaagent' flag - this points to the agent path
		layer.LaunchEnvironment.Appendf("JAVA_TOOL_OPTIONS", " ",
			"-javaagent:%s -Djma.check_interval=%s -Djma.thresholds.heap=%s", filepath.Join(layer.Path, agentPath), "5000ms", "60")

		return layer, nil
	})
}

func (w javaMemoryAssistant) Name() string {
	return w.LayerContributor.LayerName()
}
