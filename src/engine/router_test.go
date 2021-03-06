package engine

import (
	"fmt"
	"testing"
)

func Test_router(t *testing.T) {
	router := createRouter()
	router.AddRoute("GET", "/wxm/abc", tempHandler)
	handler := router.GetRoute("GET", "/wxm/abc")
	if handler != nil {
		handler(nil)
	}
}

func tempHandler(c *Context) {
	fmt.Println("tempHandler invoked")
}
