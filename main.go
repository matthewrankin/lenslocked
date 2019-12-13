package main

import (
	"fmt"
	"net/http"

	"github.com/matthewrankin/lenslocked/controllers"
	"github.com/matthewrankin/lenslocked/middleware"
	"github.com/matthewrankin/lenslocked/models"

	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "docker"
	dbname   = "lenslocked_dev"
)

func main() {

	// Create a DB connection string and then use it to create our model
	// services.
	dbInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	services, err := models.NewServices(dbInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User)
	galleriesC := controllers.NewGalleries(services.Gallery)

	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	newGallery := requireUserMw.Apply(galleriesC.New)
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.HandleFunc("/galleries", createGallery).Methods("POST")
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
