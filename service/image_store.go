package service

import (
	"bytes"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func saveImage(imageType string, imageData bytes.Buffer) (string, error) {
	imagePath := fmt.Sprintf("%s/%s", "/app/images", getImageURL(imageType))

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file: %w", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file: %w", err)
	}

	return imagePath, nil
}

func getImageURL(imageType string) string {
	if len(imageType) == 0 {
		imageType = "jpg"
	}
	return fmt.Sprintf("%s.%s", uuid.New(), imageType)
}
