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
		fmt.Println("Connect success:")
	} else {
		fmt.Println("Connect fail:", err)
	}
}
