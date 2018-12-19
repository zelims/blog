package controllers

import (
	"github.com/revel/revel"
	"github.com/zelims/blog/app/models"
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