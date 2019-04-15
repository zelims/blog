package models

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/github"
	"github.com/zelims/blog/app/database"
	"golang.org/x/oauth2"
	"html/template"
	"log"
	"net/http"
	"sort"
)

const _GITHUB_USER = "zelims" // github username for repos

var githubClient *github.Client
var httpClient *http.Client
var githubContext context.Context

type repoData []*github.Repository

func GithubAuthentication() {
	githubContext = context.Background()
	row := database.Handle.QueryRow("SELECT `github` FROM tokens")
	var token string
	err := row.Scan(&token)
	if err != nil {
		log.Printf("Could not get access token: %s", err.Error())
	}

	ts := oauth2.StaticTokenSource( &oauth2.Token{ AccessToken: token }, )
	httpClient = oauth2.NewClient(githubContext, ts)
	githubClient = github.NewClient(httpClient)
}

func GithubCalendar(username string) template.HTML {
	response, err := http.Get("https://github.com/" + username)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	out, err := document.Find(".js-yearly-contributions").Html()
	if err != nil {
		log.Printf("[!] Error getting Github calendar for %s (%s)", username, err.Error())
	}
	return template.HTML(out)
}

func GithubData() []*github.Repository {
	if githubClient == nil || githubContext == nil {
		GithubAuthentication()
	}
	repos, _, _ := githubClient.Repositories.List(githubContext, "zelims", nil)
	githubRepos := new(repoData)
	*githubRepos = repos

	sort.Sort(*githubRepos)
	return *githubRepos
}

func (p repoData) Len() int {
	return len(p)
}

func (p repoData) Less(i, j int) bool {
	return p[i].UpdatedAt.After(p[j].UpdatedAt.Time)
}

func (p repoData) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}