package controllers

import (
	"github.com/revel/revel"
	"github.com/zelims/blog/app/models"
)

type Posts struct {
	*revel.Controller
	Manage
}

func (p Posts) View() revel.Result {
	p.ViewArgs["posts"] = models.AllPosts()
	return p.checkAuth("Manage/Posts/view.html")
}
func (p Posts) New() revel.Result {
	return p.checkAuth("Manage/Posts/new.html")
}

func (p Posts) Edit(id int) revel.Result {
	return p.checkAuth("Manage/Posts/edit.html")
}