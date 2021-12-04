package controller

import (
	"Go-Echo-System/model"
	"Go-Echo-System/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type paramUserRegister struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	//Phone    string `json:"phone" validate:"required,numeric"`
	//Email    string `json:"email" validate:"required,email"`
}

type responseUserRegister struct {
	ID string `json:"_id"`
}

func UserRegister(c echo.Context) error {
	var param paramUserRegister
	if err := c.Bind(&param); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(param); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	_, found, err := model.GetUserWithUsername(param.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	if found {
		return c.String(http.StatusBadRequest, "username already exists")
	}

	user := model.User{Username: param.Username}
	user.PasswordHash, err = user.HashPassword(param.Password)
	if err != nil {
		log.Println("Hash err", err.Error())
		return c.String(http.StatusBadRequest, err.Error())
	}

	id, err := model.AddUser(user)
	if err != nil {
		log.Println("add user err:", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, responseUserRegister{
		ID: id,
	})
}

type loginValidator struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type loginResponse struct {
}

func Login(c echo.Context) error {
	var validator loginValidator
	if err := c.Bind(&validator); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(validator); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	user, found, err := model.GetUserWithUsername(validator.Username)
	if !found {
		return utils.ErrorResponse(c, http.StatusBadRequest, "user not found")
	}
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	//校验密码是否正确
	ok := user.CheckPassword(validator.Password)
	if !ok {
		return utils.ErrorResponse(c, http.StatusForbidden, "password is not correct")
	}
	//密码正确, 下面开始注册用户会话数据
	//以_gt_session 作为会话名字，获取一个session对象
	sess, _ := session.Get("_gt_session", c)
	//设置会话参数
	sess.Options = &sessions.Options{
		Path:   "/",       //所有页面都可以访问会话数据
		MaxAge: 86400 * 7, //会话有效期，单位秒
	}
	//记录会话数据, sess.Values 是map类型，可以记录多个会话数据
	sess.Values["id"] = user.ID.Hex()
	sess.Values["username"] = user.Username
	//保存用户会话数据
	sess.Save(c.Request(), c.Response())
	return c.String(http.StatusOK, "登录成功!")
}

type userInfoResponse struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
}

func GetUserInfo(c echo.Context) error {
	sess, _ := session.Get("_gt_session", c)

	//通过sess.Values读取会话数据
	id := sess.Values["id"]
	username := sess.Values["username"]

	_, found, _ := model.GetUserWithUsername(username.(string))
	if !found {
		utils.ErrorResponse(c, http.StatusBadRequest, "no such user")
	}
	return utils.SuccessResponse(c, http.StatusOK, userInfoResponse{ID: id.(string), Username: username.(string)})
}

//func UpdateUser(c echo.Context) error {
//	u := new(user)
//	if err := c.Bind(u); err != nil {
//		return err
//	}
//	id, _ := strconv.Atoi(c.Param("id"))
//	users[id].Name = u.Name
//	return c.JSON(http.StatusOK, users[id])
//}
//
//func DeleteUser(c echo.Context) error {
//	id, _ := strconv.Atoi(c.Param("id"))
//	delete(users, id)
//	return c.NoContent(http.StatusNoContent)
//}

func GetAllUserInfos(ctx echo.Context) error {

	users, err := model.GetAllUsers()
	if err != nil {
		return utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return utils.SuccessResponse(ctx, http.StatusOK, users)
}
