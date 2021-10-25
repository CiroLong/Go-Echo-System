package router

import "github.com/labstack/echo/v4"

func InitApiRouter(e *echo.Group) {
	userGroup := e.Group("/user")
	initUserGroup(userGroup)
}
