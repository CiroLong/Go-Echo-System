package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Test(c echo.Context) error {
	return c.String(http.StatusOK, "user test ok")
}
