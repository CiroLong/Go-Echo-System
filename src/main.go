package main

import (
	"Go-Echo-System/router"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	router.InitApiRouter(e.Group("/api/v1"))

	e.Logger.Fatal(e.Start(":8080"))
}
