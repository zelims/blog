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
	analyticData := models.GetCountryAnalytics()
	c.ViewArgs["unique_visitors"] = models.GetUniqueVisitors()
	c.ViewArgs["visitors"] = analyticData
	return c.checkAuth("Manage/index.html")
}