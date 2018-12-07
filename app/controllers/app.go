package controllers

import (
	"database/sql"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"log"
	"math"
	"strconv"
	"time"
)

type App struct {
	*revel.Controller
}

type Pagination struct {
	Pages		int

}

func (c App) Index() revel.Result {
	posts, size := models.Posts(1)

	pagen := &Pagination{int(math.Ceil(float64(size) / 8)) }
	pageNum := 1
	return c.Render(posts, pagen, pageNum)
}

func (c App) PagePosts() revel.Result {
	pageNum := 1
	if c.Params.Get("num") != "" {
		pageNum, _ = strconv.Atoi(c.Params.Get("num"))
	}
	posts, _ := models.Posts(pageNum)

	c.ViewArgs["posts"] = posts
	return c.RenderTemplate("Post/list.html")
}

func (c App) About() revel.Result {
	row := app.DB.QueryRow("SELECT about FROM config")
	about := ""
	if err := row.Scan(&about); err != nil || err == sql.ErrNoRows {
		c.Log.Error("Couldn't get about data: %s", err.Error())
	}
	return c.Render(about)
}

func (c App) Projects() revel.Result {
	var repos []models.RepositoryData
	if err := cache.Get("repos", &repos); err != nil {
		repos = models.Repositories()
		if repos == nil {
			hadError := "Failed to fetch GitHub data"
			return c.Render(repos, hadError)
		}
		go func() {
			err = cache.Set("repos", repos, 10*time.Minute)
			if err != nil {
				log.Printf("Error Caching Repos: %s", err.Error())
			}
		}()
	}
	return c.Render(repos)
}