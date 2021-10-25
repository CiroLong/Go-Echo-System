package router

import (
	"Go-Echo-System/controller"
	"github.com/labstack/echo/v4"
)

func initUserGroup(group *echo.Group) {
	group.GET("/test", controller.Test)
	group.POST("/login", controller.Login)
}
