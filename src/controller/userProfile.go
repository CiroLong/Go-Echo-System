package controller

import (
	"Go-Echo-System/model"
	"Go-Echo-System/utils"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type paramUserProfileChanger struct {
	Method   string `json:"method" form:"method"`
	Name     string `json:"name" form:"name"`
	Bio      string `json:"bio" form:"bio"`
	Company  string `json:"company" form:"company"`
	Location string `json:"location" form:"location"`
	Blog     string `json:"blog" form:"blog"`
}

func ChangeUserProfile(ctx echo.Context) error {
	userName := ctx.Param("username")

	user := ctx.Get("user").(model.User)

	if userName != user.Username {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "你无法修改他人的个人信息")
	}
	//binding
	var paramer paramUserProfileChanger
	if err := ctx.Bind(&paramer); err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	if paramer.Method != "put" {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "请设置method为put")
	}
	//更改
	err := user.Update(bson.M{
		"name":     paramer.Name,
		"bio":      paramer.Bio,
		"blog":     paramer.Blog,
		"company":  paramer.Company,
		"location": paramer.Location,
	})
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	//返回响应
	return utils.SuccessResponse(ctx, http.StatusOK, "更改成功")
}

type responseUserProfile struct {
	Name     string `json:"name" form:"name"`
	Bio      string `json:"bio" form:"bio"`
	Company  string `json:"company" form:"company"`
	Location string `json:"location" form:"location"`
	Blog     string `json:"blog" form:"blog"`

	Image string `json:"image"`
}

func GetUserProfile(ctx echo.Context) error {
	user := ctx.Get("user").(model.User)

	return utils.SuccessResponse(ctx, http.StatusOK, responseUserProfile{
		Name:     user.Name,
		Bio:      user.Bio,
		Company:  user.Company,
		Location: user.Location,
		Blog:     user.Blog,
		Image:    user.Image,
	})
}
