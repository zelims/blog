package controllers

import (
	"fmt"
	"github.com/revel/revel"
	"github.com/zelims/blog/app/database"
	"github.com/zelims/blog/app/models"
	"math"
)

type Post struct {
	*revel.Controller
}

func (c Post) Show(url string) revel.Result {
	post := models.GetPostByURL(url)
	if post == nil {
		return c.NotFound("Post /%s does not exist", url)
	}
	return c.Render(post)
}

func (c Post) Keywords(tag string) revel.Result {
	post := models.GetPostsByTag(tag)
	if post == nil {
		return c.NotFound("Could not find any posts tagged %s", tag)
	}
	c.ViewArgs["posts"] = post
	c.ViewArgs["search"] = "Searching posts tagged #" + tag
	c.ViewArgs["pagen"] = &models.Pagination{Pages: int(math.Ceil(float64(len(post)) / 8)) }
	c.ViewArgs["pageNum"] = 1
	return c.RenderTemplate("Post/search.html")
}

func (c Post) Search() revel.Result {
	searchInp := c.Params.Get("postSearch")

	searchQuery := "%" + searchInp + "%"
	query, err := database.Handle.Query("SELECT * FROM `posts` WHERE UPPER(content) " +
		"LIKE UPPER(?) OR UPPER(title) LIKE UPPER(?) OR UPPER(description) LIKE UPPER(?) OR FIND_IN_SET(?, `tags`)",
		searchQuery, searchQuery, searchQuery, searchInp)
	if err != nil {
		c.Log.Error("Query Error: ", err.Error())
		return c.NotFound("Could not find any posts containing %s", searchInp)
	}
	posts := models.GetPostData(query)
	search := fmt.Sprintf("Found %d posts containing %s", len(posts), searchInp)
	c.ViewArgs["pagen"] = &models.Pagination{int(math.Ceil(float64(len(posts)) / 8)) }
	c.ViewArgs["pageNum"] = 1
	return c.Render(search, posts)
}