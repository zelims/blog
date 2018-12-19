package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/revel/revel"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"github.com/zelims/blog/app/routes"
	"log"
	"strings"
)

type Sessions struct {
	*revel.Controller
}

func (c Sessions) Index() revel.Result {
	if c.currentUser() != nil {
		return c.Redirect(routes.Manage.Index())
	}
	return c.RenderTemplate("Sessions/login.html")
}

func (c Sessions) currentUser() *models.User{
	if c.Session["user"] == nil {
		return nil
	}
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	} else {
		if username, ok := c.Session["user"]; ok{
			return c.user(username.(string))
		}
	}

	return nil
}

func (c Sessions) user(username string) *models.User {
	var user models.User
	query := app.DB.QueryRow("SELECT ID,username FROM `users` WHERE username=?", username)
	err := query.Scan(&user.ID, &user.Username)
	if err != nil {
		c.Logout()
		return nil
	}
	return &user
}

func (c Sessions) userLogin(username, password string) *models.User{
	var user models.User
	query := app.DB.QueryRow("SELECT ID,username FROM `users` WHERE username=? AND password=?", username, password)
	err := query.Scan(&user.ID, &user.Username)
	if err != nil {
		log.Printf("Could not scan to user: %s", err.Error())
		return nil
	}
	return &user
}

func (c Sessions) Login(username string, password string, rememberMe bool) revel.Result {
	user := c.userLogin(username, encryptPwd(password))
	if user != nil {
		c.Session["user"] = username
		if rememberMe {
			c.Session.SetDefaultExpiration()
		} else {
			c.Session.SetNoExpiration()
		}
		c.Flash.Success("welcome " + strings.Title(username))
		return c.Redirect(routes.Manage.Index())
	}
	c.Flash.Out["username"] = username

	c.Flash.Error("Login Failed")
	return c.Redirect(routes.Sessions.Index())
}
func (c Sessions) Logout() revel.Result{
	for k:= range c.Session{
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

func encryptPwd(pwd string) string {
	h := sha256.Sum256([]byte(pwd))
	return base64.StdEncoding.EncodeToString(h[:])
}