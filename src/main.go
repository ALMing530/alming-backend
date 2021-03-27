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
	server.DELETE("/post/:id", handler.DeletePost)

	server.GET("/words", handler.GetWords)
	server.POST("/word", handler.AddWord)
	server.DELETE("/word/:id", handler.DeleteWord)
	server.GET("/translate/:word", handler.Translate)

	server.GET("/sysInfo", handler.GetInfo)

	server.GET("/timer", handler.Timer)

	server.Websocket("/:id")
	server.SetWsOnMessage(func(ws *engine.Ws, msgType int, data []byte) {
		log.Println("Custom ws callback:", string(data))
	})
	server.Run()
}
func handleHello(c *engine.Context) {
	c.Response.Write([]byte("hello"))
}
