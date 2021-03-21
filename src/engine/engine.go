package engine

import (
	"fmt"
	"net/http"
	"strings"
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
	if true {
		allowCORS(w, r)
	}
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

func (e *Engine) POST(path string, handler HandleFunc) {
	e.router.AddRoute("POST", path, handler)
}

func allowCORS(w http.ResponseWriter, r *http.Request) {
	ref := strings.TrimSuffix(r.Referer(), "/")
	w.Header().Set("Access-Control-Allow-Origin", ref)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
