package controllers

import (
	"github.com/revel/revel"
	"github.com/zelims/blog/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	posts := models.GetPosts()
	return c.Render(posts)
}

func (c App) About() revel.Result {
	return c.Render()
}

func (c App) Projects() revel.Result {
	repo := models.Repositories()
	return c.Render(repo)
}