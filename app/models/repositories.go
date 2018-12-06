package models

import (
	"encoding/json"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"
)

const _GITHUB_USER = "zelims" // github username for repos

// struct for storing repo data
type RepositoryData struct {
	Name 			string			`json:"name"`				// name of repository
	Description 	string			`json:"description"`		// description of repo
	LastUpdate 		string			`json:"updated_at"`			// total number of commits

	Language		string			`json:"language"`


	Forks			int				`json:"forks_count"`		// number of forks
	OpenIssues		int				`json:"open_issues_count"`	// current open issues
	Watchers 		int				`json:"watchers_count"`		// number of watchers
	Stars 			int				`json:"stargazers_count"`	// number of stars

	README			bool			`content:"README"`			// bool to set whether repo has README
	README_DATA    	template.HTML	`content:"README_DATA"`		// only set if above = true, contains README
}

type repoData []RepositoryData

/*****************************************************
 * GetAll (repositories)
 * @desc gets all public repositories from user
 * @return []RepositoryData: array of repoData
 *****************************************************/
func Repositories() []RepositoryData {
	// uses GitHub API to get JSON data of repos
	res, err := http.Get("https://api.github.com/users/" + _GITHUB_USER + "/repos")
	if err != nil {
		log.Printf("GH Request Error: %s", err.Error())
		return nil
	}

	// gets all the data in a raw string
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error Reading Output: %s", err.Error())
		return nil
	}

	rData := new(repoData)

	// creates an array of repoData
	err = json.Unmarshal(body, &rData) // sets the resp of repos to the above array
	if err != nil {
		fmt.Println("unmarshal error:", err)
	}

	// loop through each repo for formatting
	for m, _ := range *rData {
		// update time
		t, _ := time.Parse(time.RFC3339, (*rData)[m].LastUpdate)
		(*rData)[m].LastUpdate = time.Unix(t.Unix(), 0).Format("02 Jan 2006")

		// check if README exists
		data, _ := http.Get("https://raw.githubusercontent.com/" + _GITHUB_USER + "/" + (*rData)[m].Name + "/master/README.md")
		if data.StatusCode == 404 {
			(*rData)[m].README = false
		} else {
			(*rData)[m].README = true
			body, _ := ioutil.ReadAll(data.Body) // reads the readme and puts it in readable form.
			(*rData)[m].README_DATA = template.HTML(string(blackfriday.MarkdownCommon([]byte(body))))
		}
	}
	sort.Sort(*rData)
	return *rData // return the pointer of the struct map
}

func (p repoData) Len() int {
	return len(p)
}

func (p repoData) Less(i, j int) bool {
	t1, _ := time.Parse("02 Jan 2006", p[i].LastUpdate)
	t2, _ := time.Parse("02 Jan 2006", p[j].LastUpdate)
	return t1.After(t2)
}

func (p repoData) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}