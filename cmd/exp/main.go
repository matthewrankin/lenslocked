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
	us, err := models.NewUserService(dbInfo)
	panicOn(err)
	defer us.Close()
	us.DestructiveReset()

	// Create a user
	user := models.User{
		Name:  "Michael Scott",
		Email: "michael@dundermifflin.com",
	}
	err = us.Create(&user)
	panicOn(err)

	// Get user 1.
	foundUser, err := us.ByID(1)
	panicOn(err)
	fmt.Println(foundUser)

	// Update a user.
	user.Name = "Updated Name"
	err = us.Update(&user)
	panicOn(err)
	foundUser, err = us.ByEmail("michael@dundermifflin.com")
	panicOn(err)
	fmt.Println(foundUser)

	// Delete a user.
	err = us.Delete(foundUser.ID)
	panicOn(err)
	_, err = us.ByID(foundUser.ID)
	if err != models.ErrNotFound {
		panic("user was not deleted!")
	}

}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
