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
	"fmt"
	"regexp"
	"strings"
)

var (
	dockerImageNameRegex = regexp.MustCompile("(?:([^\\/]+)\\/)?(?:([^\\/]+)\\/)?([^@:\\/]+)(?:[@:](.+))")
)

func ParseImageName(imageName string) (string, string, string, error) {
	matches := dockerImageNameRegex.FindStringSubmatch(imageName)

	if len(matches) != 5 {
		return "", "", "", fmt.Errorf("Expected 5 matches in regex, but found %d", len(matches))
	}

	hostname := matches[1]
	namespace := matches[2]
	image := matches[3]
	tag := matches[4]

	if namespace == "" && hostname != "" {
		if !strings.Contains(hostname, ".") && !strings.Contains(hostname, ":") {
			namespace = hostname
			hostname = ""
		}
	}

	if hostname == "" {
		hostname = "index.docker.io"
	}

	if namespace == "" {
		namespace = "library"
	}

	return hostname, fmt.Sprintf("%s/%s", namespace, image), tag, nil
}
