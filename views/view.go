package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// Globals to help glob.
var (
	LayoutDir   = "views/layouts/"
	TemplateExt = ".gohtml"
)

// NewView creates a new View from the given template files.
func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// View models a view for our MVC applicaton.
type View struct {
	Template *template.Template
	Layout   string
}

// Render renders the layout template for the view.
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	return files
}
