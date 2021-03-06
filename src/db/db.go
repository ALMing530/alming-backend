package db

import (
	"database/sql"
	"fmt"
	"time"

	//init sql
	_ "github.com/go-sql-driver/mysql"
)

//DB db oprate instance
var DB *sql.DB

var sqlType = "mysql"
var dataSource = "root:123456@tcp(localhost)/alming"

func init() {
	DB, _ = sql.Open(sqlType, dataSource)
	DB.SetConnMaxLifetime(time.Minute * 3)
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err == nil {
		fmt.Println("connect success:")
	} else {
		fmt.Println("connect fail:", err)
	}
	execQuery()
}
func execQuery() {
	rows, err := DB.Query("select * from user where username=? and password=?", "wxm", "123")
	if err != nil {
		fmt.Println("An error occerred when exec query sql", err)
	}
	var name, pass string
	for rows.Next() {
		rows.Scan(&name, &pass)
		fmt.Println(name, pass)
	}
}
