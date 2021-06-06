package main

import (
	"alming_backend/src/engine"
	"alming_backend/src/handler"
	"alming_backend/src/service"
	"log"
)

func main() {
	server := engine.CreateEngine()
	api := server.Group("/api")
	api.GET("/userinfo", func(c *engine.Context) {
		c.WriteText("user:alming")
	})
	server.GET("/hello", handleHello)
	server.GET("/posts", handler.GetPosts)
	server.GET("/post/:id", handler.GetPost)
	server.POST("/upload/post", handler.MarkDownUpload)
	server.POST("/post", handler.AddPost)
	server.PUT("/post", handler.UpdatePost)
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
	server.SetWsOnOpen(func(ws *engine.Ws) {
		service.PushUsageData(ws)
	})
	server.SetWsOnClose(func(ws *engine.Ws, code int, text string) {
		service.StopPush(ws)
	})

	server.Run()
}
func handleHello(c *engine.Context) {
	c.Response.Write([]byte("hello"))
}
