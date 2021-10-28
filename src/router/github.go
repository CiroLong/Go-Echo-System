package router

import (
	"Go-Echo-System/controller"
	"github.com/labstack/echo/v4"
)

func initGithub(group *echo.Group) {
	group.POST("/session", controller.Session)
}
