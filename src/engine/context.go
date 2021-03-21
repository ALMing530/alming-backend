package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

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

//GetParam query get param with given key
func (c *Context) GetParam(key string) string {
	return c.Request.URL.Query().Get(key)
}
func (c *Context) GetParamToInt(key string) (res int, err error) {
	param := c.Request.URL.Query().Get(key)
	paramConv, err := strconv.Atoi(param)
	if err != nil {
		return -1, err
	}
	return paramConv, nil
}
func (c *Context) ParseJSONParam(obj interface{}) {
	postBody := c.PostBody()
	err := json.Unmarshal(postBody, obj)
	if err != nil {
		fmt.Println("parse json fail check your json format", err)
	}
}

//PostParam get post param with given key
func (c *Context) PostParam(key string) string {
	return c.Request.PostFormValue(key)
}

func (c *Context) PostBody() []byte {
	body := c.Request.Body
	buf := make([]byte, 128)
	var postBody []byte
	for {
		readLen, err := body.Read(buf)
		if readLen == 0 && err == io.EOF {
			break
		}
		postBody = append(postBody, buf[0:readLen]...)
	}
	return postBody
}

//WriteJSON write the json data to response
func (c *Context) WriteJSON(content interface{}) {
	c.Response.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(content)
	if err != nil {
		fmt.Println("parse json error")
	}
	_, err = c.Response.Write(bytes)
	if err != nil {
		fmt.Println("write json error")
	}
}

//WriteText write text data to response
func (c *Context) WriteText(content string) {
	c.Response.Header().Set("Content-Type", "text/plain")
	_, err := c.Response.Write([]byte(content))
	if err != nil {
		fmt.Println("write text error")
	}
}

//Param get post or get reqquest param whith given param,
func (c *Context) Param(key string) string {
	return c.Request.FormValue(key)
}

//File get file and file info
func (c *Context) File(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}
