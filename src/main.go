package main

import (
	"Go-Echo-System/router"
	"github.com/labstack/echo/middleware"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	router.InitRouter(e.Group("/api/v1"))

	e.Logger.Fatal(e.Start(":8080"))
}
