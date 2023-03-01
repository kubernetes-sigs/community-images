package cli

import (
	"fmt"

	"github.com/dims/community-images/pkg/community_images"
)

func headerLine(host string) string {
	return fmt.Sprintf("\u001B[36mImages being used in Kubernetes cluster @ %s\n", host)
}

func imageWithTag(image community_images.RunningImage) string {
	repo, img, tag, err := community_images.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	location := ""
	location = fmt.Sprintf("\u001B[36m(Location: \u001B[33m%s > %s", image.Namespace, image.Pod)
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
