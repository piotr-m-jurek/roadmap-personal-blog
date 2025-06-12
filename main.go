package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const PORT = 8008
const contentDir = "content/*.md"

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}

}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type Entry struct {
	ID      string
	Title   string
	Content string
}

func loadEntries() map[string]Entry {
	entries := map[string]Entry{}
	files, err := filepath.Glob(contentDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		parts := strings.Split(string(content), "\n\n")
		ID := strings.TrimSuffix(filepath.Base(file), ".md")

		entries[ID] = Entry{ID: ID, Title: parts[0], Content: parts[1]}
	}
	return entries
}

type Data struct {
	Page    string
	Entry   *Entry
	Entries map[string]Entry
}

func newData(page string, entry *Entry) Data {
	return Data{
		Page:    page,
		Entry:   entry,
		Entries: loadEntries(),
	}
}
func main() {

	e := echo.New()
	e.Renderer = newTemplate()
	e.Use(middleware.Logger())
	e.Static("/static", "static")

	data := newData("index", nil)
	fmt.Println(len(data.Entries))

	e.GET("/", func(c echo.Context) error {
		data = newData("index", nil)

		if c.Request().Header.Get("HX-Request") == "true" {
			return c.Render(200, "index", &data)
		}

		return c.Render(200, "page", &data)
	})

	e.GET("/blog", func(c echo.Context) error {
		data = newData("blog", nil)
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.Render(200, "blog", &data)
		}

		return c.Render(200, "page", &data)
	})

	e.GET("/blog/:id", func(c echo.Context) error {
		id := c.Param("id")
		entry, ok := data.Entries[id]

		if !ok {
			return c.Render(404, "page", &data)
		}

		data = newData("blog-entry", &entry)
		if c.Request().Header.Get("HX-Request") == "true" {
			return c.Render(200, "entry", &entry)
		}

		return c.Render(200, "page", &data)
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
