package engine

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsHelper *WsHelper

type WsOpenCallback func(ws *Ws)

type WsCloseCallback func(ws *Ws, code int, text string)

type WsMessageCallback func(ws *Ws, msgType int, data []byte)

type IDgenerater func(pathVar map[string]string, c *Context) string

type Ws struct {
	Path   string
	ConnID string
	Conn   *websocket.Conn
}

type WsHelper struct {
	Conns      map[string]*Ws
	OnOpen     WsOpenCallback
	OnMessage  WsMessageCallback
	OnClose    WsCloseCallback
	GenerageID IDgenerater
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func processWebSocket(pathVar map[string]string, c *Context, wsHelper *WsHelper) {
	conn, err := upgrader.Upgrade(c.Response, c.Request, c.Response.Header())
	if err != nil {
		log.Println("Upgrade to websocket error")
	}
	ws := &Ws{
		Path:   c.Path,
		ConnID: wsHelper.GenerageID(pathVar, c),
		Conn:   conn,
	}
	wsHelper.OnOpen(ws)
	conn.SetCloseHandler(func(code int, text string) error {
		wsHelper.OnClose(ws, code, text)
		delete(wsHelper.Conns, ws.ConnID)
		return nil
	})
	wsHelper.Conns[ws.ConnID] = ws
	go ReadMessage(ws, wsHelper.OnMessage)
}

func checkOrigin(r *http.Request) bool {
	return true
}

func ReadMessage(ws *Ws, callback WsMessageCallback) {
	for {
		msgType, msg, err := ws.Conn.ReadMessage()
		if err != nil {
			// if closeErr, ok := err.(*websocket.CloseError); ok {
			// 	log.Println(closeErr.Code, closeErr.Text)
			// }
			ws.Conn.Close()
			break
		}
		callback(ws, msgType, msg)
	}
}

func (ws WsHelper) WriteByID(connId string, data []byte) {
	w := ws.Conns[connId]
	if w != nil {
		w.Conn.WriteMessage(websocket.TextMessage, data)
	}
}

func (ws WsHelper) WriteBroadcast(data []byte) {
	for _, item := range ws.Conns {
		item.Conn.WriteMessage(websocket.TextMessage, data)
	}
}

func defaultIDgenerater(pathVar map[string]string, c *Context) string {
	return pathVar["id"]
}

func defaultWsOpenCallback(ws *Ws) {
	log.Println("Default Open handler")
}

func defaultWsCloseCallback(ws *Ws, code int, text string) {
	log.Println("Default Close handler")
}
