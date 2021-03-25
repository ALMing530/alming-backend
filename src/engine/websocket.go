package engine

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WsMessageCallback func(ws *Ws, msgType int, data []byte)
type IDgenerater func(pathVar map[string]string, c *Context) string
type Ws struct {
	Path   string
	ConnID string
	Conn   *websocket.Conn
}

type WsHelper struct {
	Conns      map[string]Ws
	OnMessage  WsMessageCallback
	GenerageID IDgenerater
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// var wsHelper WsHelper

func processWebSocket(pathVar map[string]string, c *Context, wsHelper *WsHelper) {
	conn, err := upgrader.Upgrade(c.Response, c.Request, c.Response.Header())
	if err != nil {
		log.Fatal("Upgrade to websocket error")
	}
	ws := &Ws{
		Path:   c.Path,
		ConnID: wsHelper.GenerageID(pathVar, c),
		Conn:   conn,
	}
	go ReadMessage(ws, wsHelper.OnMessage)
}
func checkOrigin(r *http.Request) bool {
	return true
}
func ReadMessage(ws *Ws, callback WsMessageCallback) {
	for {
		msgType, msg, err := ws.Conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, 1000) {
				log.Println("Socket closed")
			}
			ws.Conn.Close()
			break
		}
		callback(ws, msgType, msg)
	}
}

func defaultIDgenerater(pathVar map[string]string, c *Context) string {
	return pathVar["id"]
}
