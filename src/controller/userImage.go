package controller

import (
	"Go-Echo-System/model"
	"Go-Echo-System/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const fileNameLength = 15

//ok
func UploadIamge(ctx echo.Context) error {
	user := ctx.Get("user").(model.User)
	userName := ctx.Param("username")
	if userName != user.Username {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "你无法修改他人的个人信息")
	}
	//-----------
	// Read file
	//-----------

	// Source
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Println("Form file err:", err.Error())
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	src, err := file.Open()
	if err != nil {
		log.Println("file open err:", err.Error())
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	// Add model
	slc := strings.Split(file.Filename, ".")
	if len(slc) != 2 {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "Illegal file name")
	}
	image := model.Image{UserId: user.ID, Filename: utils.RandomString(fileNameLength) + "." + slc[1]}
	_, err = model.AddImage(image)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusInternalServerError, "add image fail, try again")
	}

	// Destination
	dst, err := os.Create("../static/images/" + image.Filename)
	if err != nil {
		log.Println("Create file err:", err.Error())
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		log.Println("io err:", err.Error())
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(ctx, http.StatusCreated, utils.H{
		"msg": fmt.Sprintf("File %s uploaded successfully.", image.Filename),
	})
}

func GetYourImage(ctx echo.Context) error {
	userName := ctx.Param("username")
	user, found, err := model.GetUserWithUsername(userName)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	if !found {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "no such user")
	}
	image, found, err := model.GetImageWithUserId(user.ID)
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}
	if !found {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, "the user has no image in the store")
	}

	return ctx.File("../static/images/" + image.Filename)
}
