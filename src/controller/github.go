package controller

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type githubLoginForm struct {
	Username string `json:"login" form:"login"`
	Password string `json:"password" form:"password"`

	RediectTo string `json:"return_to" form:"return_to"`
	Commit    string `json:"commit" form:"commit"`
}

func Session(c echo.Context) error {
	var validator githubLoginForm
	if err := c.Bind(&validator); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(validator); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
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

	sess, _ := session.Get("_gt_sess", c)
	sess.Options = &sessions.Options{
		Domain:   "github.com",
		Path:     "/",                 //所有页面都可以访问会话数据
		MaxAge:   int(time.Hour * 24), //会话有效期，单位秒
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, //Lax
	}
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusNotModified, "")
}
