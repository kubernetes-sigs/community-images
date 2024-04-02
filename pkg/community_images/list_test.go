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
)

func TestIsNamespaceExcluded(t *testing.T) {
	tests := []struct {
		name       string
		namespace  string
		exclusions []string
		expected   bool
	}{
		{
			name:       "no exclusions",
			namespace:  "foo",
			exclusions: []string{},
			expected:   false,
		},
		{
			name:       "exact match",
			namespace:  "foo",
			exclusions: []string{"foo"},
			expected:   true,
		},
		{
			name:       "exact match in list",
			namespace:  "foo",
			exclusions: []string{"one", "foo", "two"},
			expected:   true,
		},
		{
			name:       "not in list",
			namespace:  "foo",
			exclusions: []string{"one", "two"},
			expected:   false,
		},
		{
			name:       "wildcard match",
			namespace:  "foo_one",
			exclusions: []string{"foo_*"},
			expected:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := isNamespaceExcluded(test.namespace, test.exclusions)
			assert.Equal(t, test.expected, actual)
		})
	}
}
