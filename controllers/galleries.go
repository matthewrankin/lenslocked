package controllers

import (
	"fmt"
	"net/http"

	"github.com/matthewrankin/lenslocked/context"
	"github.com/matthewrankin/lenslocked/models"
	"github.com/matthewrankin/lenslocked/views"
)

// Galleries models the galleries.
type Galleries struct {
	New *views.View
	gs  models.GalleryService
}

// NewGalleries creates new galleries given the GalleryService.
func NewGalleries(gs models.GalleryService) *Galleries {
	return &Galleries{
		New: views.NewView("bootstrap", "galleries/new"),
		gs:  gs,
	}
}

// Create handles the POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	fmt.Fprintln(w, gallery)
}

// GalleryForm models the form for a gallery.
type GalleryForm struct {
	Title string `schema:"title"`
}
