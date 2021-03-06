package db

import (
	"fmt"
	"testing"
)

func Test_db(t *testing.T) {
	rows, err := DB.Query("select * from user")
	if err != nil {
		fmt.Println("An error occerred when exec query sql", err)
	}
	var id int
	var name, pass string
	var params []interface{}
	params = append(params, &id)
	params = append(params, &name)
	params = append(params, &pass)
	for rows.Next() {
		rows.Scan(params...)
		fmt.Println(id, name, pass)
	}
}
