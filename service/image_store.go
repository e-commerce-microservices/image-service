package service

import (
	"bytes"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func saveImage(imageType string, imageData bytes.Buffer) (string, error) {
	imageURL := getImageURL(imageType)
	imagePath := fmt.Sprintf("%s/%s", "/app/images", imageURL)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file: %w", err)
	}

	return imageURL, nil
}

func getImageURL(imageType string) string {
	if len(imageType) == 0 {
		imageType = "jpg"
	}
	return fmt.Sprintf("%s.%s", uuid.New(), imageType)
}
