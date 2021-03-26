package engine

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

//Engine The engine of the web service
type Engine struct {
	*router
}

//CreateEngine Create a new Engine
func CreateEngine() *Engine {
	if wsHelper == nil {
		wsHelper = &WsHelper{
			Conns:      make(map[string]*Ws),
			OnOpen:     defaultWsOpenCallback,
			OnMessage:  defaultWsOnMsg,
			OnClose:    defaultWsCloseCallback,
			GenerageID: defaultIDgenerater,
		}
	}
	return &Engine{
		router: createRouter(),
	}
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := createContext(w, r)
	if true {
		allowCORS(w, r)
	}
	if websocket.IsWebSocketUpgrade(r) {
		e.router.handleWebsocket(ctx, wsHelper)
	} else {
		e.router.handle(ctx)
	}
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
func (e *Engine) DELETE(path string, handler HandleFunc) {
	e.router.AddRoute("DELETE", path, handler)
}

func (e *Engine) Websocket(path string) {
	e.router.AddRoute("GET", path, nil)
}
func (e *Engine) SetWsOnMessage(callback WsMessageCallback) {
	wsHelper.OnMessage = callback
}
func (e *Engine) GetWsHelper() *WsHelper {
	return wsHelper
}
func allowCORS(w http.ResponseWriter, r *http.Request) {
	ref := strings.TrimSuffix(r.Referer(), "/")
	w.Header().Set("Access-Control-Allow-Origin", ref)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
func defaultWsOnMsg(ws *Ws, msgType int, data []byte) {
	fmt.Println(string(data))
}
