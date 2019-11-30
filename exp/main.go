package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	db, err := gorm.Open("postgres", dbInfo)
	panicOn(err)
	defer db.Close()
	db.LogMode(true)
	db.AutoMigrate(&User{}, &Order{})

	// Get the first user.
	var user User
	db.First(&user)
	panicOn(db.Error)

	// Create a few fake orders.
	createOrder(db, user, 1001, "Fake Description #1")
	createOrder(db, user, 9999, "Fake Description #2")
	createOrder(db, user, 8800, "Fake Description #3")
}

// User models a user.
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}

// Order models an order placed by a user.
type Order struct {
	gorm.Model
	UserID      uint
	Amount      int
	Description string
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserID:      user.ID,
		Amount:      amount,
		Description: desc,
	})
	panicOn(db.Error)
}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("What is your email?")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)
	return name, email
}

func panicOn(err error) {
	if err != nil {
		panic(err)
	}
}
