package controllers

import (
	"github.com/revel/revel"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"html/template"
	"log"
	"strconv"
	"strings"
)

type Post struct {
	*revel.Controller
}

func (c Post) Show(id int) revel.Result {
	id, err := strconv.Atoi(c.Params.Get("id"))
	if err != nil {
		return c.NotFound("Could not find %d (%s)", id, err.Error())
	}
	post := c.getPostByID(id)
	if post == nil {
		return c.NotFound("Post %d does not exist", id)
	}
	return c.Render(post)
}

func (c Post) Keywords(tag string) revel.Result {
	post := c.getPostsByTag(tag)
	if post == nil {
		return c.NotFound("Could not find any posts tagged %s", tag)
	}
	return c.Render(post, tag)

}




func (c Post) getPostsByTag(tag string) []*models.Post {
	posts := make([]*models.Post, 0)
	query, err := app.DB.Query("SELECT * FROM `posts` WHERE FIND_IN_SET(?, `tags`)", tag)
	if err != nil {
		c.Log.Error("Query Error: %s", err.Error())
		return nil
	}
	for query.Next() {
		curPost := &models.Post{}
		contentStr := ""
		if err = query.Scan(&curPost.ID, &curPost.Author, &curPost.Title, &contentStr,
			&curPost.Tags, &curPost.Date); err != nil {
			log.Printf("[!] Error scanning to post: %s", err.Error())
		}

		curPost.Tags = strings.ToLower(strings.Replace(curPost.Tags, ",", " ", -1))
		curPost.TagArr = strings.Split(curPost.Tags, " ") // creates the keyword array

		curPost.Description = template.HTML(contentStr[:500])
		curPost.Content = template.HTML(contentStr)
		curPost.FormatDate()
		posts = append(posts, curPost)
	}
	return posts
}

func (c Post) getPostByID(id int) *models.Post {
	query, err := app.DB.Query("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		c.Log.Error("Query Error: %s", err.Error())
		return nil
	}
	curPost := &models.Post{}
	for query.Next() {
		contentStr := ""
		if err = query.Scan(&curPost.ID, &curPost.Author, &curPost.Title, &contentStr,
			&curPost.Tags, &curPost.Date); err != nil {
			log.Printf("[!] Error scanning to post: %s", err.Error())
		}
		curPost.Content = template.HTML(contentStr)
		curPost.FormatDate()
	}
	return curPost
}