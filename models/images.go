package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ImageService provides the interface for the image service.
type ImageService interface {
	Create(galleryID uint, r io.Reader, filename string) error
	ByGalleryID(galleryID uint) ([]string, error)
}

// NewImageService returns a new image service.
func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) Create(galleryID uint, r io.Reader, filename string) error {
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}
	dst, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

// ByGalleryID returns all the image paths for the given gallery ID.
func (is *imageService) ByGalleryID(galleryID uint) ([]string, error) {
	path := is.imagePath(galleryID)
	strings, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}
	return strings, nil
}

func (is *imageService) imagePath(galleryID uint) string {
	return filepath.Join("images", "galleries", fmt.Sprintf("%v", galleryID))
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := filepath.Join("images", "galleries", fmt.Sprintf("%v", galleryID))
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
