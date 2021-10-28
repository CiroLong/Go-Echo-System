package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func Session(ctx echo.Context) error {
	return ctx.Redirect(http.StatusNotModified, "")
}
