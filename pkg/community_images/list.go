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
	"context"
	"fmt"
	"strings"

	"github.com/minio/pkg/wildcard"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

type RunningImage struct {
	Namespace     string
	Pod           string
	InitContainer *string
	Container     *string
	Image         string
	PullableImage string
}

func ListImages(configFlags *genericclioptions.ConfigFlags, imageNameCh chan string, ignoreNs []string, includeNs []string) ([]RunningImage, error) {
	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clientset")
	}
	var namespaces []string
	if len(includeNs) > 0 {
		namespaces = includeNs
	} else {
		clusterNamespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list namespaces")
		}
		namespaces = convertNSlistToStrList(clusterNamespaces)
	}
	runningImages := []RunningImage{}
	for _, namespace := range namespaces {
		if isNamespaceExcluded(namespace, ignoreNs) {
			continue
		}

		if imageNameCh != nil {
			imageNameCh <- fmt.Sprintf("%s/", namespace)
		}

		pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrap(err, "failed to list pods")
		}

		for _, pod := range pods.Items {
			for _, initContainerStatus := range pod.Status.InitContainerStatuses {
				pullable := initContainerStatus.ImageID
				if strings.HasPrefix(pullable, "docker-pullable://") {
					pullable = strings.TrimPrefix(pullable, "docker-pullable://")
				}
				runningImage := RunningImage{
					Pod:           pod.Name,
					Namespace:     pod.Namespace,
					InitContainer: &initContainerStatus.Name,
					Image:         initContainerStatus.Image,
					PullableImage: pullable,
				}

				if imageNameCh != nil {
					imageNameCh <- fmt.Sprintf("%s/%s", namespace, runningImage.Image)
				}
				runningImages = append(runningImages, runningImage)
			}

			for _, containerStatus := range pod.Status.ContainerStatuses {
				pullable := containerStatus.ImageID
				if strings.HasPrefix(pullable, "docker-pullable://") {
					pullable = strings.TrimPrefix(pullable, "docker-pullable://")
				}
				runningImage := RunningImage{
					Pod:           pod.Name,
					Namespace:     pod.Namespace,
					Container:     &containerStatus.Name,
					Image:         containerStatus.Image,
					PullableImage: pullable,
				}

				if imageNameCh != nil {
					imageNameCh <- fmt.Sprintf("%s/%s", namespace, runningImage.Image)
				}
				runningImages = append(runningImages, runningImage)
			}
		}
	}

	// Remove exact duplicates
	cleanedImages := []RunningImage{}
	for _, runningImage := range runningImages {
		for _, cleanedImage := range cleanedImages {
			if cleanedImage.PullableImage == runningImage.PullableImage &&
				cleanedImage.Image == runningImage.Image {
				goto NextImage
			}
		}

		cleanedImages = append(cleanedImages, runningImage)
	NextImage:
	}

	return cleanedImages, nil
}

func convertNSlistToStrList(namespaces *v1.NamespaceList) []string {
	namespacesNames := make([]string, len(namespaces.Items))
	for i, namespace := range namespaces.Items {
		namespacesNames[i] = namespace.Name
	}
	return namespacesNames
}

func isNamespaceExcluded(namespace string, excluded []string) bool {
	for _, ex := range excluded {
		if wildcard.Match(ex, namespace) {
			return true
		}
	}

	return false
}
