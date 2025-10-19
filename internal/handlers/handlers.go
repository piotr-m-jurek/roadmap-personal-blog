package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/piotr-m-jurek/roadmap-personal-blog/internal/models"
	"github.com/piotr-m-jurek/roadmap-personal-blog/internal/store"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	store   *store.Store
	session *sessions.CookieStore
}

func New(s *store.Store, session *sessions.CookieStore) *Handlers {
	return &Handlers{store: s, session: session}
}

func (h *Handlers) newData(page string, entry *models.Entry) (models.Data, error) {
	entries, err := h.store.GetEntries()
	if err != nil {
		return models.Data{}, err
	}
	return models.Data{
		Page:    page,
		Entry:   entry,
		Entries: entries,
	}, nil
}

func (h *Handlers) Index(c echo.Context) error {
	data, err := h.newData("index", nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "index", &data)
	}
	return c.Render(http.StatusOK, "page", &data)
}

func (h *Handlers) BlogIndex(c echo.Context) error {
	data, err := h.newData("blog", nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "blog", &data)
	}
	return c.Render(http.StatusOK, "page", &data)
}

func (h *Handlers) BlogEntry(c echo.Context) error {
	id := c.Param("id")
	entries, err := h.store.GetEntries()
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}

	entry, ok := entries[id]
	if !ok {
		data, _ := h.newData("404", nil) // we can ignore error here
		return c.Render(http.StatusNotFound, "page", &data)
	}

	data, err := h.newData("blog-entry", &entry)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "entry", &entry)
	}
	return c.Render(http.StatusOK, "page", &data)
}

func (h *Handlers) AdminIndex(c echo.Context) error {
	data, err := h.newData("admin/index", nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "main-content", &data)
	}
	return c.Render(http.StatusOK, "admin", &data)
}

func (h *Handlers) BlogNew(c echo.Context) error {
	data, err := h.newData("admin/new", nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}
	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "new-post", &data)
	}
	return c.Render(http.StatusOK, "admin", &data)
}

func (h *Handlers) AdminBlogNew(c echo.Context) error {
	entries, err := h.store.GetEntries()
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}

	identifier := c.FormValue("identifier")
	title := c.FormValue("title")
	content := c.FormValue("content")

	_, ok := entries[identifier]
	data, _ := h.newData("admin/new", nil)

	if title == "" || content == "" || ok || identifier == "" {
		return c.Render(http.StatusConflict, "new-post", &data)
	}

	entry := models.Entry{Title: title, Content: content, ID: identifier}

	if err := h.store.SaveEntry(entry); err != nil {
		return c.String(http.StatusInternalServerError, "could not save entry")
	}

	c.Response().Header().Set("HX-Redirect", "/admin")
	return c.Redirect(http.StatusFound, "/admin")
}

func (h *Handlers) AdminBlogEdit(c echo.Context) error {
	id := c.Param("id")
	entries, err := h.store.GetEntries()
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}

	entry, ok := entries[id]
	if !ok {
		data, _ := h.newData("404", nil)
		return c.Render(http.StatusNotFound, "page", &data)
	}

	data, err := h.newData("blog-edit", &entry)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "edit-post", &entry)
	}
	return c.Render(http.StatusOK, "admin", &data)
}

func (h *Handlers) AdminBlogUpdate(c echo.Context) error {
	id := c.Param("id")
	entries, err := h.store.GetEntries()
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}

	if _, ok := entries[id]; !ok {
		data, _ := h.newData("404", nil)
		return c.Render(http.StatusNotFound, "page", &data)
	}

	entry := models.Entry{
		Title:   c.FormValue("title"),
		Content: c.FormValue("content"),
		ID:      id,
	}

	if err := h.store.SaveEntry(entry); err != nil {
		return c.String(http.StatusInternalServerError, "could not save entry")
	}

	c.Response().Header().Set("HX-Push-Url", "/admin")

	data, err := h.newData("admin/index", nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}
	return c.Render(http.StatusOK, "main-content", &data)
}

func (h *Handlers) AdminBlogDelete(c echo.Context) error {
	id := c.Param("id")
	if err := h.store.DeleteEntry(id); err != nil {
		return c.String(http.StatusInternalServerError, "could not delete entry")
	}

	data, err := h.newData("admin/index", nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get entries")
	}
	return c.Render(http.StatusOK, "main-content", &data)
}
