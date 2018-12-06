package models

import (
	"fmt"
	"github.com/zelims/blog/app"
	"log"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	ID					int
	Author				string
	Title				string
	Description			string
	Content				string
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
		if err = query.Scan(&curPost.ID, &curPost.Author, &curPost.Title, &curPost.Content,
			&curPost.Description, &curPost.Tags, &curPost.Date); err != nil {
			log.Printf("[!] Error scanning post: %s", err.Error())
		}
		curPost.Format()
		allPosts = append(allPosts, curPost)
	}
	return allPosts
}

func (p *Post) Format() {
	p.formatDate()
	p.formatTags()
	p.formatContent()
}
func (p *Post) formatDate() {
	// starts a conversion string for UNIX --> RFC1123
	convStr, _ := strconv.ParseInt(p.Date, 10, 64)
	// setting the data from the convStr to the proper format
	p.Date = time.Unix(convStr, 0).Format("2 Jan 2006 at 3:04pm MST") //time.RFC1123
}

func (p *Post) formatTags() {
	p.Tags = strings.ToLower(strings.Replace(p.Tags, ",", " ", -1))
	p.TagArr = strings.Split(p.Tags, " ") // creates the keyword array

}

func (p *Post) formatContent() {
	p.Description += fmt.Sprintf("... <a href=\"/post/%d\" class=\"read-more\">Read more</a>", p.ID)
}