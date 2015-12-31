package main

import (
	"net/http"
	"time"

	mdl "1bwiki/model"
	"1bwiki/tmpl/special"
	"1bwiki/tmpl/special/user"

	"github.com/labstack/echo"
	"github.com/syntaqx/echo-middleware/session"
)

func register(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	flashes := session.Flashes()
	session.Save()
	return c.HTML(http.StatusOK, special.Register(val.(*mdl.User), flashes))
}

func registerHandle(c *echo.Context) error {
	if c.Form("password") == c.Form("passwordConfirm") {
		u := mdl.User{
			Name:         c.Form("username"),
			Password:     c.Form("password"),
			Registration: time.Now().Unix(),
		}
		err := mdl.CreateUser(&u)
		if err != nil {
			return logger.Error("failed to create user", "err", err)
		}
		session := session.Default(c)
		session.Set("user", u)
		session.Save()
	} else {
		session := session.Default(c)
		session.AddFlash("Passwords don't match")
		session.Save()
		return c.Redirect(http.StatusSeeOther, "/special/register")
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func login(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	return c.HTML(http.StatusOK, special.Login(val.(*mdl.User)))
}

func loginHandle(c *echo.Context) error {
	u, err := mdl.GetUserByName(c.Form("username"))
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized) // The user is doesn't exist
	}

	if u.ValidatePassword(c.Form("password")) {
		session := session.Default(c)
		session.Set("user", u)
		session.Save()
		return c.Redirect(http.StatusSeeOther, "/")
	}

	c.Response().Header().Set("Method", "GET")
	return echo.NewHTTPError(http.StatusUnauthorized)
}

func logout(c *echo.Context) error {
	session := session.Default(c)
	session.Set("user", nil)
	session.Save()
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func prefs(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*mdl.User)
	if ok {
		return c.HTML(http.StatusOK, user.Prefs(u))
	}
	return echo.NewHTTPError(http.StatusUnauthorized)
}

func prefsPasword(c *echo.Context) error {
	session := session.Default(c)
	val := session.Get("user")
	u, ok := val.(*mdl.User)
	if ok {
		return c.HTML(http.StatusOK, user.Password(u))
	}
	return echo.NewHTTPError(http.StatusUnauthorized)
}

func handlePrefsPassword(c *echo.Context) error {
	if c.Form("newpassword1") != c.Form("newpassword2") {
		// need to implement better
		return echo.NewHTTPError(http.StatusBadRequest, "password do not match")
	}
	session := session.Default(c)
	val := session.Get("user")
	u, _ := val.(*mdl.User)

	if u.ValidatePassword(c.Form("oldpassword")) {
		u.Password = c.Form("newpassword1")
		err := u.EncodePassword()
		if err != nil {
			logger.Error("registering user, encrypting password", "err", err)
		}
		err = mdl.UpdateUser(u)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.Redirect(http.StatusSeeOther, "/special/preferences")
	}

	c.Response().Header().Set("Method", "GET")
	return echo.NewHTTPError(http.StatusUnauthorized) // The user is invalid!
}
