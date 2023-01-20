package service

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

// DiskImageStore ...
type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
}

// NewDiskImageStore creates a new DiskImageStore
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
	}
}

// Save ...
func (store *DiskImageStore) Save(imageType string, imageData bytes.Buffer) (string, error) {
	imagePath := fmt.Sprintf("%s/%s/%s", os.Getenv("HOST_ADDR"), store.imageFolder, getImageURL(imageType))

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
