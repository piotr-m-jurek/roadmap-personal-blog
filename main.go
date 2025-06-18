package main

import (
	"crypto/subtle"
	"strconv"

	"github.com/piotr-m-jurek/roadmap-personal-blog/internal/handlers"
	"github.com/piotr-m-jurek/roadmap-personal-blog/internal/store"
	"github.com/piotr-m-jurek/roadmap-personal-blog/internal/view"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const PORT = 8008

func main() {

	e := echo.New()

	e.Renderer = view.New()
	e.Use(middleware.Logger())
	e.Static("/static", "static")

	s := store.New()
	h := handlers.New(s)

	e.GET("/", h.Index)
	e.GET("/blog", h.BlogIndex)
	e.GET("/blog/:id", h.BlogEntry)

	// Admin routes
	adminGroup := e.Group("/admin")
	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if subtle.ConstantTimeCompare([]byte(username), []byte("admin")) == 1 && subtle.ConstantTimeCompare([]byte(password), []byte("password")) == 1 {
			return true, nil
		}
		return false, nil
	}))

	adminGroup.GET("", h.AdminIndex)
	adminGroup.GET("/new", h.BlogNew)
	adminGroup.POST("/new", h.AdminBlogNew)
	adminGroup.GET("/edit/:id", h.AdminBlogEdit)
	adminGroup.POST("/edit/:id", h.AdminBlogUpdate)
	adminGroup.GET("/delete/:id", h.AdminBlogDelete)

	// Start server
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(PORT)))
}
