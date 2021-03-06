package engine

import "net/http"

//Context Http info context
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Method   string
	Path     string
	Header   http.Header
}

func createContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Response: w,
		Request:  r,
		Method:   r.Method,
		Path:     r.URL.Path,
	}
}
