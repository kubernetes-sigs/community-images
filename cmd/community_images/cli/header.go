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

package cli

import (
	"fmt"

	"github.com/kubernetes-sigs/community-images/pkg/community_images"
)

func headerLine(host string) string {
	return fmt.Sprintf("\u001B[36mImages being used in Kubernetes cluster @ %s\n", host)
}

func imageWithTag(image community_images.RunningImage) string {
	repo, img, tag, err := community_images.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	location := fmt.Sprintf("\u001B[36m(Location: \u001B[33m%s > %s", image.Namespace, image.Pod)
	if image.InitContainer != nil {
		location = fmt.Sprintf("%s > \u001B[33m%s", location, *image.InitContainer)
	}
	if image.Container != nil {
		marker := ">"
		if image.InitContainer != nil {
			marker = ","
		}
		location = fmt.Sprintf("%s %s \u001B[33m%s", location, marker, *image.Container)
	}
	location = fmt.Sprintf("%s\u001B[36m)", location)

	imageName := fmt.Sprintf("%s/%s", repo, img)
	truncatedImageName := imageName
	truncatedTagName := tag
	return fmt.Sprintf("%s:%s %s", truncatedImageName, truncatedTagName, location)
}

func imageWithTagPlain(image community_images.RunningImage) string {
	repo, img, tag, err := community_images.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	imageName := fmt.Sprintf("%s/%s", repo, img)
	return fmt.Sprintf("%s:%s", imageName, tag)
}
