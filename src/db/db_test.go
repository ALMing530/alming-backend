package db

import (
	"fmt"
	"testing"
)

func Test_db(t *testing.T) {
	var u user
	resMatched := QueryOne(&u, "select * from user")
	if resMatched {
		fmt.Println("QueryOne Res:", u)
	}
	var us []user
	Query(&us, "select * from user")
	fmt.Println("Query Res:", us)
	uu := user{
		Id:       2,
		Username: "alming_update2",
	}
	Exec(&uu, "update user set username=:username where id=:id")
}

type user struct {
	Id       int
	Username string
	Password string
	Time     string
}
