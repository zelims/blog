package controllers

import (
	"github.com/revel/revel"
	"github.com/zelims/blog/app/models"
	"github.com/zelims/blog/app/routes"
)

type Manage struct {
	*revel.Controller
	Sessions
}

func (c Manage) Index() revel.Result {
	c.ViewArgs["postCount"] = models.SizeOfAllPosts()
	c.ViewArgs["visitors"] = map[string]int {
		"AF": 16,
		"AL": 11,
		"DZ": 158,
		"US": 1395,
	}
	return c.checkAuth("Manage/index.html")
}

func (c Manage) Authenticated() bool {
	return !(c.currentUser() == nil)
}

func (c Manage) checkAuth(tmpl string) revel.Result {
	if c.currentUser() == nil {
		return c.Redirect(routes.Sessions.Index())
	}
	c.ViewArgs["username"] = c.currentUser().Username
	return c.RenderTemplate(tmpl)
}