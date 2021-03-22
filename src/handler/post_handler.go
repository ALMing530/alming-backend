package handler

import (
	"alming_backend/src/engine"
	"alming_backend/src/service"
)

func GetPosts(c *engine.Context) {
	posts := service.GetPosts(c)
	c.WriteJSON(posts)
}

func MarkDownUpload(c *engine.Context) {
	service.MarkDownUpload(c)
}

func GetPost(c *engine.Context) {
	post := service.GetPost(c)
	c.WriteJSON(post)
}
