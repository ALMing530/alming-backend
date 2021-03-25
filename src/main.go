package main

import (
	"alming_backend/src/engine"
	"alming_backend/src/handler"
	"log"
)

func main() {
	server := engine.CreateEngine()
	server.GET("/hello", handleHello)
	server.GET("/posts", handler.GetPosts)
	server.GET("/post/:id", handler.GetPost)
	server.POST("/post", handler.MarkDownUpload)

	server.GET("/words", handler.GetWords)
	server.POST("/word", handler.AddWord)
	server.GET("/translate/:word", handler.Translate)

	server.GET("/sysInfo", handler.GetInfo)

	server.Websocket("/:id")
	server.SetWsMsgCallback(func(ws *engine.Ws, msgType int, data []byte) {
		log.Println("Custom ws callback:", string(data))
	})
	server.Run()
}
func handleHello(c *engine.Context) {
	c.Response.Write([]byte("hello"))
}
