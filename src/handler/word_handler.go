package handler

import (
	"alming_backend/src/engine"
	"alming_backend/src/service"
)

func GetWords(c *engine.Context) {
	service.GetWords(c)
}
func AddWord(c *engine.Context) {
	service.AddWord(c)
}
func DeleteWord(c *engine.Context) {
	service.DeleteWord(c)
}

func Translate(c *engine.Context) {
	trans := service.Translate(c)
	c.WriteJSON(trans)
}
