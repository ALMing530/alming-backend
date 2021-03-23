package service

import (
	"alming_backend/src/engine"
	"alming_backend/src/repository"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func GetWords(c *engine.Context) {
	words := repository.GetWords()
	c.WriteJSON(words)
}

func AddWord(c *engine.Context) {
	word := new(repository.Word)
	c.ParseJSONParam(word)
	repository.InsertWord(word)
}

func Translate(c *engine.Context) []string {
	var trans []string
	word := c.PathParam("word")
	res, err := http.Get("http://dict.youdao.com/search?q=" + word)
	if err != nil {
		fmt.Println("白嫖有道失败")
	}
	body := res.Body
	content, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println("Read from reponse body error")
	}
	divReg, _ := regexp.Compile(`<div class="trans-container">(\n|.)*?</div>`)
	div := divReg.FindString(string(content))
	liReg, _ := regexp.Compile(`<li>(n\.|vt\.|adj\.|vi\.).*</li>`)
	li := liReg.FindAllString(div, -1)
	for _, item := range li {
		trans = append(trans, strings.TrimPrefix(strings.TrimSuffix(item, "</li>"), "<li>"))
	}
	return trans
}
