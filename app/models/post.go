package models

import (
	"github.com/zelims/blog/app"
	"html/template"
	"log"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	ID					int
	Author				string
	Title				string
	Description			template.HTML
	Content				template.HTML
	Tags				string
	TagArr 				[]string
	Date				string
}

func GetPosts() []*Post {
	allPosts := make([]*Post, 0)

	query, err := app.DB.Query("SELECT * FROM posts")
	if err != nil {
		log.Printf("[!] Error getting posts: %s", err.Error())
	}
	for query.Next() {
		curPost := &Post{}
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
		allPosts = append(allPosts, curPost)
	}
	return allPosts
}

func GetPost(id int) *Post {

	curPost := &Post{}
	return curPost
}
func (p *Post) FormatDate() {
	// starts a conversion string for UNIX --> RFC1123
	convStr, _ := strconv.ParseInt(p.Date, 10, 64)
	// setting the data from the convStr to the proper format
	p.Date = time.Unix(convStr, 0).Format("2 Jan 2006 at 3:04pm MST") //time.RFC1123
}