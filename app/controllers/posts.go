package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/kataras/iris/core/errors"
	"github.com/olahol/go-imageupload"
	"github.com/revel/revel"
	"github.com/zelims/blog/app/database"
	"github.com/zelims/blog/app/models"
	"github.com/zelims/blog/app/routes"
	"log"
	"net/http"
	"os"
	"strings"
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

func checkURLExists(url string) int {
	count := 0
	err := database.Handle.QueryRow("SELECT COUNT(*) FROM posts WHERE friendly_url LIKE ?", url).Scan(&count)
	if err != nil {
		log.Printf("Could not query row %s", err.Error())
		return -1
	}
	return count
}

func (p Posts) Create() revel.Result {
	if !p.Authenticated() {
		return p.Redirect(routes.Sessions.Index())
	}

	postValues, err := p.getFormValues()
	if err != nil {
		// error handle
		return p.Redirect(routes.Posts.View())
	}

	_, err = database.Handle.NamedExec(`INSERT INTO posts (ID, author, title, content, description, friendly_url, tags, banner, images, date)` +
		` VALUES (:id,:author,:title,:content,:desc,:url,:tags,:banner,:images,:date)`,
		map[string]interface{}{
			"id":          	nil,
			"author":      	p.currentUser().Username,
			"title":       	postValues["title"],
			"content":     	postValues["content"],
			"desc":        	postValues["description"],
			"url":		   	postValues["url"],
			"tags":        	postValues["tags"],
			"banner":		postValues["banner"],
			"images":		postValues["images"],
			"date":        	time.Now().Unix(),
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

	postValues, err := p.getFormValues()
	if err != nil {
		log.Printf("Could not get form values: %s", err.Error())
		p.Flash.Error(fmt.Sprintf("Couldn't modify post: %s", err.Error()))
		return p.Redirect(routes.Posts.View())
	}

	curPost, err := models.PostByID(id)
	if err != nil {
		log.Printf("Could not get current post in database: %s", err.Error())
	}

	if postValues["banner"] != "" && postValues["banner"] != curPost.Banner {
		err := os.Remove(fmt.Sprintf("public/images/posts/%s.jpg", curPost.Banner))
		if err != nil {
			log.Printf("Could not remove previous banner: %s", err.Error())
		}
	}

	query := "title=:title, content=:content, description=:desc, tags=:tags, last_update=:date"
	queryData := map[string]interface{}{
		"id":      			id,
		"title":   			postValues["title"],
		"content": 			postValues["content"],
		"desc":    			postValues["description"],
		"tags":    			postValues["tags"],
		"date":    			time.Now().Unix(),
	}

	if postValues["banner"] != "" {
		query += ",banner=:banner"
		queryData["banner"] = postValues["banner"]
	}

	_, err = database.Handle.NamedExec(`UPDATE posts SET `+query+` WHERE id=:id`, queryData)

	if err != nil {
		log.Printf("Could not update post %d: %s", id, err.Error())
		p.Flash.Error(fmt.Sprintf("Couldn't modify post: %s", err.Error()))
		return p.RenderTemplate(routes.Posts.Edit(id))
	}
	p.Flash.Success("Post successfully modified!")
	return p.Redirect(routes.Posts.Edit(id))
}

func (p Posts) getFormValues() (map[string]string, error) {
	formValues := make(map[string]string)
	var err error

	formValues["title"]				= p.Params.Form.Get("post-title")
	formValues["content"] 			= p.Params.Form.Get("post-content")
	formValues["description"] 		= p.Params.Form.Get("post-description")
	formValues["tags"] 				= p.Params.Form.Get("post-tags")

	titleLen := len(formValues["title"])
	if titleLen > 64 {
		titleLen = 64
	}
	friendlyURL := strings.ToLower(strings.Replace(formValues["title"][0:titleLen], " ", "-", -1))
	count := checkURLExists(friendlyURL + "%")
	if count > 0 {
		friendlyURL = fmt.Sprintf("%s-%d", friendlyURL, count+1)
	} else if count == -1 {
		log.Printf("Could not create post as the URL messed up")
		p.Flash.Error(fmt.Sprintf("The friendly url check failed, please report this issue!"))
		return nil, errors.New("Friendly url check failed")
	}
	formValues["url"] = friendlyURL

	formValues["banner"], err = p.getBannerImage(formValues["title"])
	if err != nil {
		log.Printf("Create error (banner): %s", err.Error())
		p.Flash.Error(fmt.Sprintf("Create error: %s", err.Error()))
		return nil, err
	}

	formValues["images"], err = p.getPostImages(formValues["title"])
	if err != nil {
		log.Printf("Create error (images): %s", err.Error())
		p.Flash.Error(fmt.Sprintf("Create error: %s", err.Error()))
		return nil, err
	}
	return formValues, err
}

func (p Posts) throwEditErr(err error, id int) bool {
	if err != nil {
		log.Printf("Could not modify post #%d: %s", id, err.Error())
		p.Flash.Error(fmt.Sprintf("Couldn't modify post: %s", err.Error()))
		return true
	}
	return false
}

func (p Posts) getBannerImage(title string) (name string, err error) {
	image, err := imageupload.Process(p.Request.In.GetRaw().(*http.Request), "post-banner")

	if err == http.ErrMissingFile {
		return "", nil
	} else if err != nil {
		return
	}

	titleHash := md5.Sum([]byte(title))
	titleHashStr := hex.EncodeToString(titleHash[:])

	fileHash := md5.Sum([]byte(image.Filename))
	fileHashStr := hex.EncodeToString(fileHash[:])

	fileName := md5.Sum([]byte(titleHashStr + fileHashStr))
	image.Filename = hex.EncodeToString(fileName[:])

	err = image.Save(fmt.Sprintf("public/images/posts/%s.jpg", image.Filename))

	if err != nil {
		log.Printf("Failed to save banner: %s", err.Error())
		return
	}

	return image.Filename, err
}

func (p Posts) getPostImages(title string) (names string, error error) {

	// return "" if there is no file uploaded
	return "", nil
}