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
	Content				string
	Description			string
	URL					string		`db:"friendly_url"`
	Tags				string
	TagArr 				[]string 	`db:"-"`
	TagsValue			string	 	`db:"-"`
	Banner				string
	Images				string
	Date				string
	Updated				*string		`db:"last_update"`
}

type FileInfo struct {
	ContentType string
	Filename    string
	RealFormat  string `json:",omitempty"`
	Resolution  string `json:",omitempty"`
	Size        int
	Status      string `json:",omitempty"`
}

func SizeOfAllPosts() int {
	size := 0
	query := app.DB.QueryRow("SELECT COUNT(*) FROM posts")
	if err := query.Scan(&size); err != nil {
		log.Printf("Failed to count rows: %s", err.Error())
		return -1
	}
	return size
}

func PostByID(id int) (post Post, err error) {
	err = app.DB.Get(&post, "SELECT * FROM posts WHERE id = ?", id)
	return
}

func AllPosts() []*Post {
	allPosts := make([]*Post, 0)
	err := app.DB.Select(&allPosts,"SELECT * FROM posts ORDER BY date DESC")
	if err != nil {
		log.Printf("AllPosts(): %s", err.Error())
	}
	formatPosts(allPosts)
	return allPosts
}

func Posts(offset int) ([]*Post, int) {
	posts := make([]*Post, 0)
	err := app.DB.Select(&posts,"SELECT * FROM posts ORDER BY date DESC LIMIT ?, 8", (offset - 1) * 8)
	if err != nil {
		log.Printf("Posts(%d): %s", offset, err.Error())
	}
	formatPosts(posts)
	return posts, SizeOfAllPosts()
}

func formatPosts(posts []*Post) {
	for _, p := range posts {
		p.Format()
	}
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

	if p.Updated != nil {
		convStr, _ = strconv.ParseInt(*p.Updated, 10, 64)
		*p.Updated = time.Unix(convStr, 0).Format("2 Jan 2006 at 3:04pm MST") //time.RFC1123
	}
}

func (p *Post) formatTags() {
	p.TagsValue = p.Tags
	p.Tags = strings.ToLower(strings.Replace(p.Tags, ",", " ", -1))
	p.TagArr = strings.Split(p.Tags, " ") // creates the keyword array

}

func (p *Post) formatContent() {
	p.Description += fmt.Sprintf("... <a href=\"/post/%s\" class=\"read-more\">Read more</a>", p.URL)
}