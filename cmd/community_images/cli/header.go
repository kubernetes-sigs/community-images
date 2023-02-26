package cli

import (
	"fmt"
	"github.com/dims/community-images/pkg/community_images"
)

func headerLine() string {
	return fmt.Sprintf("Image")
}

func imageWithTag(image community_images.RunningImage) string {
	repo, img, tag, err := community_images.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	imageName := fmt.Sprintf("%s/%s", repo, img)
	truncatedImageName := imageName
	truncatedTagName := tag
	return fmt.Sprintf("%s:%s", truncatedImageName, truncatedTagName)
}
