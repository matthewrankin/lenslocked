package views

import "html/template"

// NewView creates a new View from the given template files.
func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/footer.gohtml",
		"views/layouts/bootstrap.gohtml")
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
