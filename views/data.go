package views

import (
	"log"

	"github.com/matthewrankin/lenslocked/models"
)

// Alert levels allowed.
const (
	AlertLvlError   = "danger"
	AlertLvlWarning = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric is displayed when any random error is encountered by our
	// backend.
	AlertMsgGeneric = "Something went wrong. Please try again, and contact us if the problem persists."
)

// PublicError is the interface for public errors.
type PublicError interface {
	error
	Public() string
}

// Alert is used to render Bootstrap Alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data // to come in.
type Data struct {
	Alert *Alert
	User  *models.User
	Yield interface{}
}

// SetAlert sets an alert on the Data type.
func (d *Data) SetAlert(err error) {
	var msg string
	if pErr, ok := err.(PublicError); ok {
		msg = pErr.Public()
	} else {
		log.Println(err)
		msg = AlertMsgGeneric
	}
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

// AlertError sets a cusotm error message.
func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}
