package utils

import "github.com/labstack/echo/v4"

type ResponseData struct {
	Success bool        `json:"success"`
	Msg     string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ErrorResponse(c echo.Context, code int, msg string) error {
	return c.JSON(code, ResponseData{
		Success: false,
		Msg:     msg,
		Data:    nil,
	})
}

func SuccessResponse(c echo.Context, code int, data interface{}) error {
	return c.JSON(code, ResponseData{
		Success: true,
		Msg:     "",
		Data:    data,
	})
}
