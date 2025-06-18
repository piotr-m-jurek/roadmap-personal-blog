package view

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	tmpl *template.Template
}

func New() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}
