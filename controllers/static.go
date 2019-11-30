package controllers

import (
	"github.com/matthewrankin/lenslocked/views"
)

// NewStatic creates the static views.
func NewStatic() *Static {
	return &Static{
		Home:    views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

// Static models the various static views.
type Static struct {
	Home    *views.View
	Contact *views.View
}
