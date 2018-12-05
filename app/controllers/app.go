package controllers

import (
	"database/sql"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"time"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	posts := models.GetPosts()
	return c.Render(posts)
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
		go cache.Set("repos", repos, 10*time.Minute)
	}
	return c.Render(repos)
}