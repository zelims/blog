package controllers

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/revel/revel"
	"github.com/zelims/blog/app/database"
	"github.com/zelims/blog/app/models"
	"github.com/zelims/blog/app/routes"
	"log"
	"net/http"
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
	query := database.Handle.QueryRow("SELECT ID,username FROM `users` WHERE username=?", username)
	err := query.Scan(&user.ID, &user.Username)
	if err != nil {
		c.Logout()
		return nil
	}
	return &user
}

func (c Sessions) userLogin(username, password string) *models.User{
	var user models.User
	query := database.Handle.QueryRow("SELECT ID,username FROM `users` WHERE username=? AND password=?", username, password)
	err := query.Scan(&user.ID, &user.Username)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Could not scan to user: %s", err.Error())
		}
		return nil
	}
	return &user
}

func (c Sessions) Edit() revel.Result {
	var profile models.UserProfile
	err := database.Handle.Get(&profile, "SELECT * FROM config")
	if err != nil {
		log.Printf("Couldn't get about data: %s", err.Error())
	}
	c.ViewArgs["profile"] = profile
	return c.checkAuth("Sessions/edit.html")
}

func (c Sessions) SaveProfile() revel.Result {
	if !c.Authenticated() {
		return c.Redirect(routes.Sessions.Index())
	}
	_, err := database.Handle.NamedExec(`UPDATE config SET name=:name,location=:loc,about=:about,github=:gh,twitter=:tw,`+
			`instagram=:ig,linkedin=:li`,
		map[string]interface{}{
			"name": 		c.Params.Form.Get("user-name"),
			"loc":	 		c.Params.Form.Get("user-location"),
			"about": 		c.Params.Form.Get("user-about"),
			"gh": 			c.Params.Form.Get("user-github"),
			"tw":	 		c.Params.Form.Get("user-twitter"),
			"ig": 			c.Params.Form.Get("user-instagram"),
			"li": 			c.Params.Form.Get("user-linkedin"),
		})
	if err != nil {
		c.Flash.Error(fmt.Sprintf("Could not update profile - %s", err.Error()))
		return c.Redirect(routes.Sessions.Edit())
	}
	c.Flash.Success("Updated profile!")
	return c.Redirect(routes.Sessions.Edit())
}

func (c Sessions) Login(username string, password string, rememberMe bool) revel.Result {
	if !models.ValidateRecaptcha(c.Request.In.GetRaw().(*http.Request)) {
		c.Flash.Error("Recaptcha Failed")
		return c.Redirect(routes.Sessions.Index())
	}
	user := c.userLogin(username, encryptPwd(password))
	if user != nil {
		c.Session["user"] = username
		if rememberMe {
			c.Session.SetNoExpiration()
		} else {
			c.Session.SetDefaultExpiration()
		}
		c.Flash.Success("Welcome " + username + "!")
		return c.Redirect(routes.Manage.Index())
	}
	c.Flash.Out["username"] = username
	if rememberMe {
		c.Flash.Out["rememberMe"] = "checked"
	}
	c.Flash.Error("Login Failed")
	return c.Redirect(routes.Sessions.Index())
}
func (c Sessions) Logout() revel.Result{
	for k:= range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

func encryptPwd(pwd string) string {
	h := sha256.Sum256([]byte(pwd))
	return base64.StdEncoding.EncodeToString(h[:])
}

func (c Sessions) Authenticated() bool {
	return !(c.currentUser() == nil)
}

func (c Sessions) checkAuth(tmpl string) revel.Result {
	if c.currentUser() == nil {
		return c.Redirect(routes.Sessions.Index())
	}
	c.ViewArgs["username"] = c.currentUser().Username
	return c.RenderTemplate(tmpl)
}