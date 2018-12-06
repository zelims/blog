package controllers

import (
	"database/sql"
	"github.com/revel/revel"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"log"
	"strconv"
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
	query, err := app.DB.Query("SELECT * FROM `posts` WHERE FIND_IN_SET(?, `tags`)", tag)
	if err != nil {
		c.Log.Error("Query Error: %s", err.Error())
		return nil
	}

	return c.getPostData(query)
}

func (c Post) getPostByID(id int) *models.Post {
	query, err := app.DB.Query("SELECT * FROM posts WHERE id = ?", id)
	if err != nil {
		c.Log.Error("Query Error: %s", err.Error())
		return nil
	}
	curPost := &models.Post{}
	for query.Next() {
		if err = query.Scan(&curPost.ID, &curPost.Author, &curPost.Title, &curPost.Content,
			&curPost.Description, &curPost.Tags, &curPost.Date); err != nil {
			log.Printf("[!] Error scanning to post: %s", err.Error())
		}
		curPost.Format()
	}
	return curPost
}

func (c Post) getPostData(query *sql.Rows) []*models.Post {
	posts := make([]*models.Post, 0)
	for query.Next() {
		curPost := &models.Post{}
		if err := query.Scan(&curPost.ID, &curPost.Author, &curPost.Title, &curPost.Content,
			&curPost.Description, &curPost.Tags, &curPost.Date); err != nil {
			log.Printf("[!] Error scanning to post: %s", err.Error())
		}
		curPost.Format()
		posts = append(posts, curPost)
	}
	return posts
}