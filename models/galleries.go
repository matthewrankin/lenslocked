package models

import "github.com/jinzhu/gorm"

// Gallery models a gallery resource.
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

// GalleryService provides the interface the gallery service.
type GalleryService interface {
	GalleryDB
}

// GalleryDB provides the interface for interacting with the database for a
// gallery.
type GalleryDB interface {
	Create(gallery *Gallery) error
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	// TODO
	return nil
}
