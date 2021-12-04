package middleware

import (
	"Go-Echo-System/model"
	"Go-Echo-System/utils"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func UserValidator(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		sess, _ := session.Get("_gt_session", ctx)
		//通过sess.Values读取会话数据
		id := sess.Values["id"]
		if id == nil {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "no cookie")
		}

		user, found, _ := model.GetUserWithID(id.(string))
		if !found {
			return utils.ErrorResponse(ctx, http.StatusBadRequest, "your cookie is wrong")
		}
		ctx.Set("user", user)
		if err := next(ctx); err != nil {
			ctx.Error(err)
		}
		return nil
	}
}
