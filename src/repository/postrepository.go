package repository

import (
	"alming_backend/src/db"
	"bytes"
	"io"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"gitlab.com/golang-commonmark/markdown"
)

type Post struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Summary    string `json:"summary"`
	Oringin    string `json:"oringin"`
	Format     string `json:"format"`
	PostType   int    `json:"postType"`
	Visits     int    `json:"visits"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"update_time"`
	WordCount  int    `json:"wordCount"`
}

func CreateMdPostDefault(title string, oringin []byte) (post Post) {
	b := bytes.NewBuffer(oringin)
	if title == "" {
		line, err := b.ReadString('\n')
		if err == nil {
			reg, _ := regexp.Compile(`#*\s`)
			title = reg.ReplaceAllString(line, "")
		}
	}
	summary := ""
	for {
		line, err := b.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		if !(strings.HasPrefix(line, "#") || strings.HasPrefix(line, "```")) {
			if utf8.RuneCountInString(summary) < 100 {
				summary += line
			} else {
				summary = string([]rune(summary)[:96]) + "..."
				break
			}
		}
	}
	// unsafe := blackfriday.Run(oringin, blackfriday.WithRenderer(&blackfriday.HTMLRenderer{}))
	// p := bluemonday.UGCPolicy()
	// p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	// html := p.SanitizeBytes(unsafe)
	html := markdown.New(markdown.HTML(true)).RenderToString(oringin)
	post.Title = title
	post.Summary = summary
	post.Oringin = string(oringin)
	mdParsed := html
	post.Format = string(mdParsed)
	post.PostType = 0
	post.Visits = 0
	post.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	post.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	post.WordCount = len(string(oringin))
	return
}
func GetPosts(posts *[]Post) {
	db.Query(posts, "select * from post")
}
func GetPost(id int) *Post {
	var sql string = `select * from post where id=?`
	var post Post
	db.QueryOne(&post, sql, id)
	return &post
}

func InsertPost(post *Post) {
	var sql string = `insert into post  values
		(0,:title,:summary,:oringin,:format,:post_type,
		:visits,:create_time,:update_time,:word_count)`
	db.Exec(post, sql)
}
