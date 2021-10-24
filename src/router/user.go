package router

import (
	"Go-Echo-System/controller"
	"github.com/labstack/echo"
)

func initUserGroup(group *echo.Group) {
	group.GET("/test", controller.Test)
}
