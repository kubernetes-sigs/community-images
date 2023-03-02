/*
Copyright 2023 The Kubernetes Authors.

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

package community_images

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseImageName(t *testing.T) {
	url, image, tag, err := ParseImageName("redis:4")
	require.NoError(t, err)
	assert.Equal(t, "index.docker.io", url)
	assert.Equal(t, "library/redis", image)
	assert.Equal(t, "4", tag)

	url, image, tag, err = ParseImageName("k8s.gcr.io/cluster-proportional-autoscaler-amd64:1.1.2-r2")
	require.NoError(t, err)
	assert.Equal(t, "k8s.gcr.io", url)
	assert.Equal(t, "library/cluster-proportional-autoscaler-amd64", image)
	assert.Equal(t, "1.1.2-r2", tag)

	url, image, tag, err = ParseImageName("quay.io/coreos/grafana-watcher:v0.0.8")
	require.NoError(t, err)
	assert.Equal(t, "quay.io", url)
	assert.Equal(t, "coreos/grafana-watcher", image)
	assert.Equal(t, "v0.0.8", tag)

	url, image, tag, err = ParseImageName("grafana/grafana:5.0.1")
	require.NoError(t, err)
	assert.Equal(t, "index.docker.io", url)
	assert.Equal(t, "grafana/grafana", image)
	assert.Equal(t, "5.0.1", tag)

	url, image, tag, err = ParseImageName("postgres:10.0")
	require.NoError(t, err)
	assert.Equal(t, "index.docker.io", url)
	assert.Equal(t, "library/postgres", image)
	assert.Equal(t, "10.0", tag)

	url, image, tag, err = ParseImageName("localhost:32000/postgres:10.0")
	require.NoError(t, err)
	assert.Equal(t, "localhost:32000", url)
	assert.Equal(t, "library/postgres", image)
	assert.Equal(t, "10.0", tag)
}
