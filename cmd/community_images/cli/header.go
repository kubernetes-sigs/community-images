package cli

import (
	"fmt"
	outdated "github.com/dims/community-images/pkg/community_images"
)

const (
	maxImageLength = 50
	maxTagLength   = 50
)

func headerLine() string {
	return fmt.Sprintf("Image")
}

func runningImage(image outdated.RunningImage) string {
	repo, img, tag, err := outdated.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	imageName := fmt.Sprintf("%s/%s", repo, img)
	truncatedImageName := imageName
	if len(truncatedImageName) > maxImageLength {
		truncatedImageName = fmt.Sprintf("%s...", truncatedImageName[0:maxImageLength-3])
	}

	truncatedTagName := tag
	if len(tag) > maxTagLength {
		truncatedTagName = fmt.Sprintf("%s...", truncatedTagName[0:maxTagLength-3])
	}

	return fmt.Sprintf("%s:%s", truncatedImageName, truncatedTagName)
}

func completedImage(image outdated.RunningImage, checkResult *outdated.CheckResult) string {
	repo, img, tag, err := outdated.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	imageName := fmt.Sprintf("%s/%s", repo, img)
	truncatedImageName := imageName
	if len(truncatedImageName) > maxImageLength {
		truncatedImageName = fmt.Sprintf("%s...", truncatedImageName[0:maxImageLength-3])
	}

	truncatedTagName := tag
	if len(tag) > maxTagLength {
		truncatedTagName = fmt.Sprintf("%s...", truncatedTagName[0:maxTagLength-3])
	}

	truncatedLatestTagName := checkResult.LatestVersion
	if len(truncatedLatestTagName) > maxTagLength {
		truncatedLatestTagName = fmt.Sprintf("%s...", truncatedLatestTagName[0:maxTagLength-3])
	}

	return fmt.Sprintf("%s:%s", truncatedImageName, truncatedTagName)

}

func erroredImage(image outdated.RunningImage, checkResult *outdated.CheckResult) string {
	repo, img, tag, err := outdated.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	imageName := fmt.Sprintf("%s/%s", repo, img)
	truncatedImageName := imageName
	if len(truncatedImageName) > maxImageLength {
		truncatedImageName = fmt.Sprintf("%s...", truncatedImageName[0:maxImageLength-3])
	}

	truncatedTagName := tag
	if len(tag) > maxTagLength {
		truncatedTagName = fmt.Sprintf("%s...", truncatedTagName[0:maxTagLength-3])
	}

	message := "Unable to get image data"
	if checkResult != nil {
		message = checkResult.CheckError
	}
	return fmt.Sprintf("%s:%s %s", truncatedImageName, truncatedTagName, message)

}
