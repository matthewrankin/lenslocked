package controllers

import (
	"fmt"
	"net/http"

	"github.com/matthewrankin/lenslocked/models"
	"github.com/matthewrankin/lenslocked/views"
)

// SignupForm models the data for the signup form.
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// NewUsers handles creating a new user.
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

// Users models a user of the web app.
type Users struct {
	NewView *views.View
	us      *models.UserService
}

// New is used to render the form where a user can create a new user account.
//
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create us used to process the signup form when a user tries to create a new
// user account.
//
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User is", user)
}
