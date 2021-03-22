package repository

import "alming_backend/src/db"

type Word struct {
	Id       int    `json:"id"`
	En       string `json:"en"`
	Cn       string `json:"cn"`
	Familiar int    `json:"familiar"`
}

//temp:All repository tend to refactor to this style.
//GetWords get all word from db
func GetWords() *[]Word {
	var sql = `select * from words`
	var words []Word
	db.Query(&words, sql)
	return &words
}
