package controllers

import (
	"github.com/revel/revel"
	"github.com/zelims/blog/app/routes"
)

type Manage struct {
	*revel.Controller
	Sessions
}

func (c Manage) Index() revel.Result {
	if c.currentUser() == nil {
		return c.Redirect(routes.Sessions.Index())
	}
	username := c.currentUser().Username
	return c.Render(username)
}