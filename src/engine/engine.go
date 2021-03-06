package engine

import (
	"fmt"
	"net/http"
)

//Engine The engine of the web service
type Engine struct {
	*router
}

//CreateEngine Create a new Engine
func CreateEngine() *Engine {
	return &Engine{
		router: createRouter(),
	}
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.router.handle(createContext(w, r))
}

//Run start the web application
func (e *Engine) Run() {
	err := http.ListenAndServe(":53000", e)
	if err != nil {
		fmt.Println("Start web server fail")
	}
}

//GET Register a get router
func (e *Engine) GET(path string, handler HandleFunc) {
	e.router.AddRoute("GET", path, handler)
}
