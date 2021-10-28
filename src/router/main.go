package router

import (
	mymiddleware "Go-Echo-System/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
func Router() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(mymiddleware.SessionMiddleWare)

	e.Validator = &CustomValidator{
		validator: validator.New(),
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	InitApiRouter(e.Group("/"))
	return e
}

func InitApiRouter(e *echo.Group) {
	userGroup := e.Group("/api/v1/user")
	initUserGroup(userGroup)
	initGithub(e)
}
