package models

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/zelims/blog/app"
	"golang.org/x/oauth2"
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
	row := app.DB.QueryRow("SELECT `github_token` FROM config")
	var token string
	err := row.Scan(&token)
	if err != nil {
		log.Printf("Could not get access token: %s", err.Error())
	}

	ts := oauth2.StaticTokenSource( &oauth2.Token{ AccessToken: token }, )
	httpClient = oauth2.NewClient(githubContext, ts)
	githubClient = github.NewClient(httpClient)
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