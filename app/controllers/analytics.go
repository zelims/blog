package controllers

import "github.com/revel/revel"

type Analytics struct {
	*revel.Controller
	Manage
}

func (a Analytics) View() revel.Result {
	return a.checkAuth("Manage/Analytics/view.html")
}

func (a Analytics) Posts() revel.Result {
	return a.checkAuth("Manage/Analytics/posts.html")
}