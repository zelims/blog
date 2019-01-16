package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"github.com/zelims/blog/app/routes"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
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
			"id":          nil,
			"author":      p.currentUser().Username,
			"title":       p.Params.Form.Get("post-title"),
			"content":     p.Params.Form.Get("post-content"),
			"desc":        p.Params.Form.Get("post-description"),
			"tags":        p.Params.Form.Get("post-tags"),
			"date":        time.Now().Unix(),
			"last_update": nil,
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

	query := "title=:title, content=:content, description=:desc, tags=:tags, last_update=:date"
	queryData := map[string]interface{}{
		"id":      id,
		"title":   p.Params.Form.Get("post-title"),
		"content": p.Params.Form.Get("post-content"),
		"desc":    p.Params.Form.Get("post-description"),
		"tags":    p.Params.Form.Get("post-tags"),
		"date":    time.Now().Unix(),
	}
	if p.Params.Form.Get("post-banner") != "" {
		request := p.Request.In.GetRaw().(*http.Request)

		file, handle, err := request.FormFile("post-banner")
		if p.throwEditErr(err, id) {
			return p.RenderTemplate(routes.Posts.Edit(id))
		}
		defer file.Close()

		if valid, fname := p.HandleImageUpload(file, handle, id); valid == true {
			query += "banner=:banner"
			queryData["banner"] = fname
			log.Printf("Updated banner: %s", queryData)
		} else {
			log.Printf("[!] Could not upload banner: %s", err.Error())
		}
	}

	_, err := app.DB.NamedExec(`UPDATE posts SET `+query+` WHERE id=:id`, queryData)

	if err != nil {
		log.Printf("Could not update post %d: %s", id, err.Error())
		p.Flash.Error(fmt.Sprintf("Couldn't modify post: %s", err.Error()))
		return p.RenderTemplate(routes.Posts.Edit(id))
	}
	p.Flash.Success("Post successfully modified!")
	return p.Redirect(routes.Posts.Edit(id))
}

func (p Posts) throwEditErr(err error, id int) bool {
	if err != nil {
		log.Printf("Could not update banner %d: %s", id, err.Error())
		p.Flash.Error(fmt.Sprintf("Couldn't modify post: %s", err.Error()))
		return true
	}
	return false
}

const (
	_      = iota
	KB int = 1 << (10 * iota)
	MB
	GB
)

type FileInfo struct {
	ContentType string
	Filename    string
	RealFormat  string `json:",omitempty"`
	Resolution  string `json:",omitempty"`
	Size        int
	Status      string `json:",omitempty"`
}

func (p Posts) HandleImageUpload(file multipart.File, handle *multipart.FileHeader, id int) (bool, string) {
	p.Validation.Required(file)
	p.Validation.MinSize(file, 2*KB).
		Message("Minimum a file size of 2KB expected")
	p.Validation.MaxSize(file, 5*MB).
		Message("File cannot be larger than 5MB")

	data, err := ioutil.ReadAll(file)
	if p.throwEditErr(err, id) {
		return false, ""
	}

	//format := handle.Header.Get("Content-Type")
	log.Printf("%s", handle.Filename)
	//p.Validation.Required(format == "jpeg" || format == "png").Key("post-banner").
	//	Message("JPEG or PNG file format is expected")

	err = ioutil.WriteFile("/public/img/posts/"+strconv.Itoa(id)+"/"+handle.Filename, data, 0666)
	if p.throwEditErr(err, id) {
		return false, ""
	}

	log.Printf("%x", handle.Filename)
	return true, handle.Filename
}