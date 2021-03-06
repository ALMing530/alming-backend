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
	UpdateTime string `json:"updateTime"`
	WordCount  int    `json:"wordCount"`
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
	var sql = `insert into post  values(0,:title,:summary,:oringin,:format,:post_type,
					:visits,:create_time,:update_time,:word_count)`
	db.Exec(post, sql)
}
func DeletePost(post *Post) {
	var sql = `delete from post where id=:id`
	db.Exec(post, sql)
}

func UpdatePost(post *Post) {
	var sql = `update post set title=:title,summary=:summary,oringin=:oringin,format
					=:format,update_time=:update_time,word_count=:word_count where id = :id`
	db.Exec(post, sql)
}

func CreateMdPostDefault(title string, oringin []byte) (post *Post) {
	b := bytes.NewBuffer(oringin)
	if title == "" {
		line, err := b.ReadString('\n')
		if err == nil {
			reg, _ := regexp.Compile(`#*\s`)
			tmp := reg.ReplaceAllString(line, "")
			if tmp == "" {
				title = "NONE"
			}
		} else {
			title = "NONE"
		}
	}
	html := markdown.New(markdown.HTML(true)).RenderToString(oringin)
	post.Title = title
	post.Summary = generateSummary(b)
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

func FillMissiongField(post *Post) {
	b := bytes.NewBuffer([]byte(post.Oringin))
	post.Summary = generateSummary(b)
	post.PostType = 1
	post.Visits = 0
	post.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	post.UpdateTime = time.Now().Format("2006-01-02 15:04:05")
	post.WordCount = len(string(post.Oringin))
}

func CreatBlankPost() *Post {
	return &Post{
		PostType:   0,
		Visits:     0,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		UpdateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
}

func generateSummary(content *bytes.Buffer) string {
	summary := ""
	for {
		line, err := content.ReadString('\n')
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
	return summary
}
