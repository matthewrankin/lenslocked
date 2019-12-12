package main

import (
	"fmt"

	"github.com/matthewrankin/lenslocked/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "lenslocked_dev"
)

func main() {
	dbInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(dbInfo)
	if err != nil {
		panic(err)
	}
	defer services.User.Close()
	services.User.DestructiveReset()
}
