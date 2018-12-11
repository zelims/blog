package controllers

import "github.com/revel/revel"

type Posts struct {
	*revel.Controller
	Manage
}

func (p Posts) View() revel.Result {
	return p.checkAuth("Manage/Posts/view.html")
}
func (p Posts) New() revel.Result {
	return p.checkAuth("Manage/Posts/new.html")
}