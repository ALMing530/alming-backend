package main

import (
	"alming_backend/src/engine"
)

func main() {
	engine := engine.CreateEngine()
	engine.GET("/hello", handleHello)
	engine.Run()
}
func handleHello(c *engine.Context) {
	c.Response.Write([]byte("hello"))
}
