package controllers

import (
	"database/sql"
	"fmt"
	"github.com/revel/revel"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"log"
	"math"
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
	c.ViewArgs["posts"] = post
	c.ViewArgs["search"] = "Searching posts tagged #" + tag
	c.ViewArgs["pagen"] = &Pagination{int(math.Ceil(float64(len(post)) / 8)) }
	c.ViewArgs["pageNum"] = 1
	return c.RenderTemplate("Post/search.html")
}

func (c Post) Search() revel.Result {
	searchInp := c.Params.Get("postSearch")
	searchQuery := "%" + searchInp + "%"
	query, err := app.DB.Query("SELECT * FROM `posts` WHERE UPPER(content) " +
		"LIKE UPPER(?) OR UPPER(title) LIKE UPPER(?) OR UPPER(description) LIKE UPPER(?) OR FIND_IN_SET(?, `tags`)",
		searchQuery, searchQuery, searchQuery, searchInp)
	if err != nil {
		c.Log.Error("Query Error: ", err.Error())
		return c.NotFound("Could not find any posts containing %s", searchInp)
	}
	posts := c.getPostData(query)
	search := fmt.Sprintf("Found %d posts containing %s", len(posts), searchInp)
	c.ViewArgs["pagen"] = &Pagination{int(math.Ceil(float64(len(posts)) / 8)) }
	c.ViewArgs["pageNum"] = 1
	return c.Render(search, posts)
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
			&curPost.Description, &curPost.Tags, &curPost.Banner, &curPost.Images, &curPost.Date,
			&curPost.Updated); err != nil {
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
			&curPost.Description, &curPost.Tags, &curPost.Banner, &curPost.Images, &curPost.Date,
			&curPost.Updated); err != nil {
			log.Printf("[!] Error scanning to post: %s", err.Error())
		}
		curPost.Format()
		posts = append(posts, curPost)
	}
	return posts
}