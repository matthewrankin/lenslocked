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
	panicOn(err)
	defer services.User.Close()
	services.User.DestructiveReset()

	// Create a user
	user := models.User{
		Name:     "Michael Scott",
		Email:    "michael@dundermifflin.com",
		Password: "bestboss",
	}
	err = services.User.Create(&user)
	panicOn(err)

	// Verify that the user has a Remember and RememberHash.
	fmt.Printf("%+v\n", user)
	if user.Remember == "" {
		panic("Invalid remember token")
	}

	// Now verify that we can lookup a user with that remember token.
	user2, err := services.User.ByRemember(user.Remember)
	panicOn(err)
	fmt.Printf("%+v\n", *user2)
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
