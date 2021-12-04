package router

import (
	"Go-Echo-System/controller"
	"Go-Echo-System/middleware"
	"github.com/labstack/echo/v4"
)

func initUserGroup(group *echo.Group) {
	group.GET("/test", controller.Test)
	group.POST("/register", controller.UserRegister)
	group.POST("/login", controller.Login)
	group.GET("/info", controller.GetUserInfo, middleware.UserValidator)

	group.GET("/all", controller.GetAllUserInfos)

	// profiles
	group.POST("/:username", controller.ChangeUserProfile, middleware.UserValidator)
	group.GET("/:username/profile", controller.GetUserProfile, middleware.UserValidator)
}
