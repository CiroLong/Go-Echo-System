package controller

import (
	"Go-Echo-System/model"
	"Go-Echo-System/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type githubLoginForm struct {
	Username string `json:"login" form:"login" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`

	RedirectTo string `json:"return_to" form:"return_to" validate:"required"`
	Token      string `json:"authenticity_token" form:"authenticity_token"`
	Commit     string `json:"commit" form:"commit" `
}

// GithubSession 总结：_gh_sess 在登陆时使用，确定登录步骤
// user_session 在登录后使用，确定用户状态
func GithubSession(c echo.Context) error {
	var validator githubLoginForm
	if err := c.Bind(&validator); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(validator); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	//!< 登录验证
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
	//!< 首先处理headers
	c.Response().Header().Add("cache-control", "no-cache")
	c.Response().Header().Add("location", validator.RedirectTo)
	c.Response().Header().Add("permissions-policy", "interest-cohort=()")
	c.Response().Header().Add("referrer-policy", "origin-when-cross-origin, strict-origin-when-cross-origin")
	c.Response().Header().Add("server", "GitHub.com")
	c.Response().Header().Add("strict-transport-security", "max-age=31536000; includeSubdomains; preload")
	c.Response().Header().Add("vary", "X-PJAX, X-PJAX-Container")
	c.Response().Header().Add("vary", "Accept-Encoding, Accept, X-Requested-With")
	c.Response().Header().Add("x-content-type-options", "nosniff")
	c.Response().Header().Add("x-frame-options", "deny")
	//c.Response().Header().Add("x-github-request-id", "8C58:6C7C:157CFC:1D344C:617CF5C0") // ?? what's up?
	c.Response().Header().Add("x-xss-protection", "0")
	//!< 然后处理set-cookie
	{
		now := time.Now()
		//user_session 要与 __Host-user_session_same_site 相同,只改变SameSite的模式
		//user_session
		id, _ := strconv.Atoi(user.ID.Hex())
		jwtToken := utils.GenerateJWT(uint(id), user.Username)
		userSession := &http.Cookie{
			Name:       "user_session",
			Value:      jwtToken,
			Domain:     "github.com",
			Path:       "/",
			Expires:    now.AddDate(0, 0, 14),
			RawExpires: now.AddDate(0, 0, 14).Format(time.UnixDate),
			MaxAge:     int(time.Hour * 24 * 14), //两周
			Secure:     true,
			HttpOnly:   true,
			SameSite:   http.SameSiteLaxMode, //Lax
		}
		c.SetCookie(userSession)
		//__Host-user_session_same_site
		hostUserSessionSameSite := &http.Cookie{
			Name:       "__Host-user_session_same_site",
			Value:      jwtToken,
			Domain:     "github.com",
			Path:       "/",
			Expires:    now.AddDate(0, 0, 14),
			RawExpires: now.AddDate(0, 0, 14).Format(time.UnixDate),
			MaxAge:     int(time.Hour * 24 * 14), //两周
			Secure:     true,
			HttpOnly:   true,
			SameSite:   http.SameSiteStrictMode, //Strict
		}
		c.SetCookie(hostUserSessionSameSite)
		// _gt_sess
		gtSess, _ := session.Get("_gt_sess", c)
		gtSess.Options = &sessions.Options{
			Domain:   "github.com",
			Path:     "/",                 //所有页面都可以访问会话数据
			MaxAge:   int(time.Hour * 24), //会话有效期，单位秒
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode, //Lax
		}
		gtSess.Values["isAdmin"] = true
		gtSess.Values["username"] = user.Username
		gtSess.Save(c.Request(), c.Response())
		//dotcom_user
		dotcomUser := &http.Cookie{
			Name:       "dotcom_user",
			Value:      validator.Username,
			Path:       "/",
			Domain:     ".github.com",
			Expires:    now.AddDate(1, 0, 0),
			RawExpires: now.AddDate(1, 0, 0).Format(time.UnixDate),
			MaxAge:     int(time.Hour * 24 * 365),
			Secure:     true,
			HttpOnly:   true,
			SameSite:   http.SameSiteLaxMode,
		}
		c.SetCookie(dotcomUser)
		//logged_in
		loggedIn := &http.Cookie{
			Name:       "logged_in",
			Value:      "yes",
			Path:       "/",
			Domain:     ".github.com",
			Expires:    now.AddDate(1, 0, 0),
			RawExpires: now.AddDate(1, 0, 0).Format(time.UnixDate),
			MaxAge:     int(time.Hour * 24 * 365),
			Secure:     true,
			HttpOnly:   true,
			SameSite:   http.SameSiteLaxMode,
		}
		c.SetCookie(loggedIn)
		//tz
		tz, err := c.Cookie("tz")
		if err == nil {
			tz.Domain = "github.com"
			tz.HttpOnly = true
			c.SetCookie(tz)
		}
		//has_recent_activity
		hasRecentActivity := &http.Cookie{
			Name:       "has_recent_activity",
			Value:      "1",
			Path:       "/",
			Domain:     "github.com",
			Expires:    now.Add(time.Hour),
			RawExpires: now.Add(time.Hour).Format(time.UnixDate),
			MaxAge:     int(time.Hour),
			Secure:     true,
			HttpOnly:   true,
			SameSite:   http.SameSiteLaxMode,
		}
		c.SetCookie(hasRecentActivity)
	}

	//return utils.SuccessResponse(c, http.StatusOK, "session ok")
	return c.Redirect(http.StatusNotModified, validator.RedirectTo) // 	重定向
}

func GithubLogin(c echo.Context) error {
	//!< headers
	c.Response().Header().Add("cache-control", "no-cache")
	c.Response().Header().Add("location", "https://github.com/")
	c.Response().Header().Add("permissions-policy", "interest-cohort=()")
	c.Response().Header().Add("referrer-policy", "origin-when-cross-origin, strict-origin-when-cross-origin")
	c.Response().Header().Add("server", "GitHub.com")
	c.Response().Header().Add("strict-transport-security", "max-age=31536000; includeSubdomains; preload")
	c.Response().Header().Add("vary", "X-PJAX, X-PJAX-Container")
	c.Response().Header().Add("vary", "Accept-Encoding, Accept, X-Requested-With")
	c.Response().Header().Add("x-content-type-options", "nosniff")
	c.Response().Header().Add("x-frame-options", "deny")
	//c.Response().Header().Add("x-github-request-id", "8C58:6C7C:157CFC:1D344C:617CF5C0") // ?? what's up?
	c.Response().Header().Add("x-xss-protection", "0")

	//!< 验证
	sess, _ := session.Get("_gh_sess", c)
	isAdmin, ok := sess.Values["isAdmin"]
	if !ok || !isAdmin.(bool) {
		return c.Redirect(http.StatusNotModified, "https://github.com/") // 	重定向
	}
	username, _ := sess.Values["username"]
	_, found, _ := model.GetUserWithUsername(username.(string))
	if !found {
		return utils.ErrorResponse(c, http.StatusForbidden, "no such user")
	}
	//!< set-cookies

	// _gt_sess
	gtSess, _ := session.Get("_gt_sess", c)
	gtSess.Options = &sessions.Options{
		Domain:   "github.com",
		Path:     "/",                 //所有页面都可以访问会话数据
		MaxAge:   int(time.Hour * 24), //会话有效期，单位秒
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, //Lax
	}
	gtSess.Values["end"] = true
	gtSess.Save(c.Request(), c.Response())
	now := time.Now()
	//has_recent_activity
	hasRecentActivity := &http.Cookie{
		Name:       "has_recent_activity",
		Value:      "1",
		Path:       "/",
		Domain:     "github.com",
		Expires:    now.Add(time.Hour),
		RawExpires: now.Add(time.Hour).Format(time.UnixDate),
		MaxAge:     int(time.Hour),
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
	}
	c.SetCookie(hasRecentActivity)

	return c.Redirect(http.StatusNotModified, "https://github.com/")
}
