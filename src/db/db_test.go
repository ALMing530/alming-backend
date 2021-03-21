package db

import (
	"fmt"
	"testing"
)

func Test_db(t *testing.T) {
	var u user
	resMatched := QueryOne(&u, "select * from user where id=1")
	if resMatched {
		fmt.Println("QueryOne Res:", u)
	}
	var us []user
	Query(&us, "select * from user")
	fmt.Println("Query Res:", us)
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
}
