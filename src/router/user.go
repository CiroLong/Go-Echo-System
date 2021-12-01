package router

import (
	"Go-Echo-System/controller"
	"github.com/labstack/echo/v4"
)

func initUserGroup(group *echo.Group) {
	group.GET("/test", controller.Test)
	group.POST("/register", controller.UserRegister)
	group.POST("/login", controller.Login)
	group.GET("/info", controller.GetUserInfo)

	group.GET("/all", controller.GetAllUserInfos)
}
