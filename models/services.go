package models

import "github.com/jinzhu/gorm"

// NewServices creates all the services using the given connection info.
func NewServices(connectionInfo string) (*Services, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &Services{
		User:    NewUserService(db),
		Gallery: &galleryGorm{},
	}, nil
}

// Services contains all the services.
type Services struct {
	Gallery GalleryService
	User    UserService
}
