package db

import (
	"fmt"
	"testing"
)

func Test_db(t *testing.T) {
	var u user
	success := QueryOne(&u, "select * from user where id=1")
	if success {
		fmt.Println("QueryOne Res:\n", u)
	}

	var us []user
	success = Query(&us, "select * from user")
	if success {
		fmt.Println("Query Res:\n", us)
	}

	var ps []Post
	var sql = `select p.id,p.title ,t.id as tid,t.name from post p,tags t,post_tag pt where p.id =pt.post_id and t.id = pt.tag_id `
	QueryOneToMany(&ps, sql, "id", "tid")
	fmt.Println("QueryOneToMany Res:\n", ps)

	uu := user{
		Id:       2,
		Username: "alming_update",
	}
	Exec(&uu, "update user set username=:username where id=:id")
}

type user struct {
	Id       int
	Nickname string
	Username string
	Password string
	Email    string
	RegTime  string
}
type Post struct {
	Id    int
	Title string
	Tags  []Tag
}
type Tag struct {
	Tid  int
	Name string
}
