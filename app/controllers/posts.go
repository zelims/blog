package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"github.com/zelims/blog/app/routes"
	"log"
	"time"
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
	post, err := models.PostByID(id)
	if err != nil  {
		log.Printf("Couldn't find Post #%d: %s", id, err.Error())
		p.Flash.Error(err.Error())
		return p.Redirect(routes.Posts.View())
	}
	p.ViewArgs["post"] = post
	return p.checkAuth("Manage/Posts/edit.html")
}

func (p Posts) Create() revel.Result {
	if !p.Authenticated() {
		return p.Redirect(routes.Sessions.Index())
	}
	_, err := app.DB.NamedExec(`INSERT INTO posts (ID, author, title, content, description, tags, date)` +
		` VALUES (:id,:author,:title,:content,:desc,:tags,:date)`,
		map[string]interface{}{
			"id": 		nil,
			"author": 	p.currentUser().Username,
			"title": 	p.Params.Form.Get("post-title"),
			"content": 	p.Params.Form.Get("post-content"),
			"desc": 	p.Params.Form.Get("post-description"),
			"tags": 	p.Params.Form.Get("post-tags"),
			"date": 	time.Now().Unix(),
		})
	if err != nil {
		log.Printf("Could not insert into posts: %s", err.Error())
		p.Flash.Error(fmt.Sprintf("Couldn't create post: %s", err.Error()))
		return p.RenderTemplate(routes.Posts.View())
	}
	p.Flash.Success("Post created!")
	return p.Redirect(routes.Posts.View())
}
func (p Posts) Modify(id int) revel.Result {
	if !p.Authenticated() {
		return p.Redirect(routes.Sessions.Index())
	}
	// do modify sql statements
	return p.Redirect(routes.Posts.Edit(id))
}