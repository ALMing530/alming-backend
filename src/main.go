package main

import (
	"alming_backend/src/engine"
	"alming_backend/src/handler"
)

func main() {
	engine := engine.CreateEngine()
	engine.GET("/hello", handleHello)
	engine.GET("/posts", handler.GetPosts)
	engine.GET("/post/:id", handler.GetPost)
	engine.POST("/post", handler.MarkDownUpload)

	engine.GET("/words", handler.GetWords)
	engine.POST("/word", handler.AddWord)
	engine.GET("/translate/:word", handler.Translate)
	engine.Run()
}
func handleHello(c *engine.Context) {
	c.Response.Write([]byte("hello"))
}
