package controllers

import (
	"github.com/google/go-github/github"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
	"github.com/zelims/blog/app"
	"github.com/zelims/blog/app/models"
	"html/template"
	"log"
	"math"
	"strconv"
	"time"
)

type App struct {
	*revel.Controller
}

func init() {
	revel.InterceptFunc(models.TrackUser, revel.AFTER, &App{})
	revel.InterceptFunc(models.TrackUser, revel.AFTER, &Post{})
}

func (c App) Index() revel.Result {
	posts, size := models.Posts(1)

	pagen := &models.Pagination{int(math.Ceil(float64(size) / 8)) }
	pageNum := 1
	c.ViewArgs["posts"] = posts
	c.ViewArgs["pagen"] = pagen
	c.ViewArgs["pageNum"] = pageNum
	return c.Render()
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
	var profile models.UserProfile
	if err := cache.Get("profile", &profile); err != nil {
		err := app.DB.Get(&profile, "SELECT * FROM config")
		if err != nil {
			log.Printf("Couldn't get about data: %s", err.Error())
		}
		go func() {
			err = cache.Set("profile", profile, 15*time.Minute)
			if err != nil {
				log.Printf("Error caching profile: %s", err.Error())
			}
		}()
	}
	c.ViewArgs["profile"] = profile
	return c.Render()
}

func (c App) Projects() revel.Result {
	return c.Render()
}

func (c App) Repositories() revel.Result {
	return c.Render()
}

func (c App) GitHub() revel.Result {
	var githubUsername string
	if err := cache.Get("githubUsername", &githubUsername); err != nil {
		err := app.DB.Get(&githubUsername, "SELECT github FROM config")
		if err != nil {
			log.Printf("[!] Could not get Github username (%s)", err.Error())
		}
		go func() {
			err = cache.Set("githubUsername", githubUsername, 10*time.Minute)
			if err != nil {
				log.Printf("Error Caching GH Username: %s", err.Error())
			}
		}()
	}

	var repos []*github.Repository
	if err := cache.Get("github", &repos); err != nil {
		repos = models.GithubData()
		go func() {
			err = cache.Set("repos", repos, 10*time.Minute)
			if err != nil {
				log.Printf("Error Caching Repos: %s", err.Error())
			}
		}()
	}

	var calendar template.HTML
	if err := cache.Get("github-calendar", &calendar); err != nil {
		calendar = models.GithubCalendar(githubUsername)
		go func() {
			err = cache.Set("github-calendar", calendar, 10*time.Minute)
			if err != nil {
				log.Printf("Error Caching Calendar: %s", err.Error())
			}
		}()
	}

	c.ViewArgs["calendar"] = calendar
	c.ViewArgs["repos"] = repos
	return c.RenderTemplate("App/ajax_data/github.html")
}