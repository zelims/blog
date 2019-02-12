package controllers

import (
	"github.com/revel/revel"
	"github.com/zelims/blog/app/models"
	"github.com/zelims/blog/app/routes"
)

type Analytics struct {
	*revel.Controller
	Manage
}

func (a Analytics) View() revel.Result {
	a.ViewArgs["visitors"] = models.GetAnalyticData()
	return a.checkAuth("Manage/Analytics/view.html")
}

func (a Analytics) Post() revel.Result {
	a.ViewArgs["postvisitors"] = models.PostsWithAnalytics()
	return a.checkAuth("Manage/Analytics/post.html")
}

func (a Analytics) Posts() revel.Result {
	return a.checkAuth("Manage/Analytics/posts.html")
}

func (a Analytics) User(uuid string) revel.Result {
	analytics := models.AnalyticsByUUID(uuid)
	if len(analytics) == 0 {
		a.Flash.Error("Could not get analytics for that user")
		return a.Redirect(routes.Analytics.View())
	}
	a.ViewArgs["user_uuid"] = uuid
	a.ViewArgs["analytics_user"] = analytics
	return a.checkAuth("Manage/Analytics/user.html")
}