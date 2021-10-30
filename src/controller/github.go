package controller

import (
	"Go-Echo-System/model"
	"Go-Echo-System/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type githubLoginForm struct {
	Username string `json:"login" form:"login" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`

	RediectTo string `json:"return_to" form:"return_to" validate:"required"`
	Token     string `json:"authenticity_token" form:"authenticity_token"`
	Commit    string `json:"commit" form:"commit" `
}

/* github 的cookie有那些字段
%2FEA99DSnWRcCZYyDowR76aur2vCOI%2FccgmhRrx3%2B%2F%2FMz8LMBwA6O4Gp4%2FD4%2BgyILrQi%2FbS3X5k0LPu%2FAGDkk2Wga2oWKRIHX9hPzPZVkMWHSymjBVfU%2Fb2u%2FogpTpWyOtS75k3rqoU1PHruKQlT7GwG%2F7D577rA0MBQvWJejVlwL5D2p%2BlL9i3Gx%2B7V4meY4R1AMkVKCXP8Io4JYUdBuu59%2F%2BTE8jpqh1b8tn1FQdhN%2FFBQMDQz9eeN%2FRndGAgEsifh7%2BG5hnmjWsMrzz1XG%2BA%3D%3D--7Si7wJF%2B8kwZNtL5--3Iwfc3RdFKiHNXCS%2B2kTNw%3D%3D	gi
Name	Value Domain   path        Expires  ...


has_recent_activity	1	github.com	/	2021-10-28T09:12:54.000Z	101	✓	✓	Lax		Medium
logged_in	yes	.github.com	/	2022-10-28T08:12:54.000Z	113	✓	✓	Lax		Medium
dotcom_user	CiroLong	.github.com	/	2022-10-28T08:12:54.000Z	120	✓	✓	Lax		Medium

__Host-user_session_same_site	GU9lnyCdKVoWfVuzPg0nGEl3OVIRY85yzOsqpN4xkhBjIcYU	github.com	/	2021-11-11T08:12:54.000Z	161	✓	✓	Strict		Medium
user_session	GU9lnyCdKVoWfVuzPg0nGEl3OVIRY85yzOsqpN4xkhBjIcYU	github.com	/	2021-11-11T08:12:54.000Z	141	✓	✓	Lax		Medium
_gh_sess	9w6fcamUqk3FAuVE5%2BD1X9a%2Ft4nL0tS5x%2BQxsNad8p…..	github.com	/	Session	1299	✓	✓	Lax		Medium
*/

// 总结：_gh_sess 在登陆时使用，确定登录步骤
// user_session 在登录后使用，确定用户状态
func Session(c echo.Context) error {
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
	c.Response().Header().Add("location", validator.RediectTo)
	c.Response().Header().Add("permissions-policy", "interest-cohort=()")
	c.Response().Header().Add("referrer-policy", "origin-when-cross-origin, strict-origin-when-cross-origin")
	c.Response().Header().Add("server", "GitHub.com")
	c.Response().Header().Add("strict-transport-security", "max-age=31536000; includeSubdomains; preload")
	c.Response().Header().Add("vary", "X-PJAX, X-PJAX-Container")
	c.Response().Header().Add("vary", "Accept-Encoding, Accept, X-Requested-With")
	c.Response().Header().Add("x-content-type-options", "nosniff")
	c.Response().Header().Add("x-frame-options", "deny")
	c.Response().Header().Add("x-github-request-id", "8C58:6C7C:157CFC:1D344C:617CF5C0") // ?? value?
	c.Response().Header().Add("x-xss-protection", "0")
	//!< 然后处理set-cookie
	{
		now := time.Now()
		//user_session 要与 __Host-user_session_same_site 相同,只改变SameSite的模式
		//user_session
		jwtToken := utils.GenerateJWT(user.ID, user.Username)
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
	return c.Redirect(http.StatusNotModified, validator.RediectTo) // 	重定向
}
