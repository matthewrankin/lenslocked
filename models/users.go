package models

import (
	"regexp"
	"strings"

	"github.com/matthewrankin/lenslocked/internal/pkg/hash"
	"github.com/matthewrankin/lenslocked/internal/pkg/rand"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Blank import needed here.
	"golang.org/x/crypto/bcrypt"
)

var (
	_ UserService = &userService{}
)

const (
	// ErrNotFound is returned when a resource cannot be found in the database.
	ErrNotFound modelError = "models: resource not found"
	// ErrIDInvalid is returned when an invalid ID is provided to a method like
	// Delete.
	ErrIDInvalid modelError = "models: ID provided was invalid"
	// ErrPasswordIncorrect is returned when an invalid password is used when
	// attempting to authenticate a user.
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	// ErrPasswordRequired is returned when a create is attempted without a user
	// password provided.
	ErrPasswordRequired modelError = "models: password is required"
	// ErrEmailRequired is returned when an email address is not provided when
	// creating a user.
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailInvalid is returned when an email address provided does not match
	// any of our requirements.
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailTaken is returned when an update or create is attempted with an
	// email address that is already in use.
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrPasswordTooShort is returned when a user tries to set a password that
	// is less than 8 characters long.
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	// ErrRememberRequired is returned when a create or update is attempted
	// without a user remember token hash
	ErrRememberRequired modelError = "models: remember token is required"
	// ErrRememberTooShort is returned when a remember token is not at least 32
	// bytes.
	ErrRememberTooShort modelError = "models: remember token must be at least 32 bytes"
	userPwPepper                   = "secret-random-string"
)

const (
	hmacSecretKey = "secret-hmac-key"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

// UserDB is used to interact with the users database.
//
// For pretty much all single user queries:
// - If the user is found, we will return a nil error
// - If the user is not found, we will return ErrNotFound
// - If there is another error, we will return an error with // more
//   information about what went wrong. This may not be // an error generated
//   by the models package.
//
// For single user queries, any error but ErrNotFound should // probably result
// in a 500 error until we make "public" facing errors.
type UserDB interface {
	// Methods for querying for single users
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)
	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

// userGorm represents our database interaction layer and implements the
// UserDB interface fully.
type userGorm struct {
	db *gorm.DB
}

// userValidator is our validation layer that validates and normalizes data
// before passing it on to the next UserDB in our interface chain.
type userValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

func newUserValidator(udb UserDB, hmac hash.HMAC) *userValidator {
	return &userValidator{
		UserDB: udb,
		hmac:   hmac,
		emailRegex: regexp.MustCompile(
			`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func (uv *userValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

// User models a user
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}

// UserService is a set of methods used to manipulate and work with the user
// model.
type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

type userService struct {
	UserDB
}

// NewUserService creates a new UserService.
func NewUserService(db *gorm.DB) UserService {
	ug := &userGorm{db}
	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
	}
}

// ByID will look up a user with the provided ID.  If the user is found, we
// will return a nil error If the user is not found, we will return
// ErrNotFound. If there is another error, we will return an error with more
// information about what went wrong. This may not be an error generated by
// the models package.  As a general rule, any error but ErrNotFound should
// probably result in a 500 error.
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail will normalize an email address before passing it on to the database
// layer to perform the query.
func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	err := runUserValFns(&user, uv.normalizeEmail)
	if err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

// ByEmail looks up a user with the given email address and returns that user.
// If the user is found, we will return a nil error If the user is not found,
// we will return ErrNotFound If there is another error, we will return an
// error with more information about what went wrong. This may not be an error
// generated by the models package.
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Update will hash a remember token if it is provided.
func (uv *userValidator) Update(user *User) error {
	err := runUserValFns(
		user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil

	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

// Update will update the provided user with all of the data in the provided
// user object.
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

// Delete will delete the user with the provided ID
func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFns(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// Delete will delete the user with the provided ID.
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Create will create the provided user and backfill data like the ID,
// CreatedAt, and UpdatedAt fields.
func (uv *userValidator) Create(user *User) error {
	err := runUserValFns(
		user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail,
	)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) idGreaterThan(n uint) userValFn {
	return userValFn(func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}
		return nil
	})
}

// Create will create the provided user and backfill data like the ID,
// CreatedAt, and UpdatedAt fields.
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// first will query using the provided gorm.DB and it will get the first item
// returned and place it into dst. If nothing is found in the query, it will
// return ErrNotFound.
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

// Authenticate can be used to authenticate a user with the provided email
// address and password. If the email address provided is invalid, this will
// return nil, ErrNotFound. If the password provided is invalid, this will
// return nil, ErrPasswordIncorrect. If the email and password are both valid,
// this will return user, nil. Otherwise if another error is encountered this
// will return nil, error.
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash), []byte(password+userPwPepper))
	switch err {
	case nil:
		return foundUser, nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return nil, ErrPasswordIncorrect
	default:
		return nil, err
	}
}

// ByRemember looks up a user with the given remember token and returns that
// user. This method will handle hashing the token for us.
func (uv *userValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFns(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

// ByRemember looks up a user with the given remember token and returns that
// user. This method will handle hashing the token for us.
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type userValFn func(*User) error

// bcryptPassword will hash a user's password with an
// app-wide pepper and bcrypt, which salts for us.
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes,
		bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func runUserValFns(user *User, fns ...userValFn) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	// If we get here that means we found a user w/ this email address, so we
	// need to see if this is the same user we are updating, or if we have a
	// conflict.
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *userValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}
