package engine

//HandleFunc The request handler
type HandleFunc func(*Context)

type router struct {
	routes  map[string]*trieNode
	handler map[string]HandleFunc
}

func createRouter() *router {
	return &router{
		routes:  make(map[string]*trieNode),
		handler: make(map[string]HandleFunc),
	}
}

func (r *router) AddRoute(method string, path string, handler HandleFunc) {
	_, ok := r.routes[method]
	if !ok {
		r.routes[method] = rootNode()
	}
	r.routes[method].AddNode(path)
	r.handler[method+":"+path] = handler
}
func (r *router) GetRoute(method string, path string) HandleFunc {
	route, ok := r.routes[method]
	if !ok {
		return nil
	}
	node := route.search(path)
	if node != nil {
		return r.handler[method+":"+node.pattern]
	}
	return nil
}
func (r *router) handle(c *Context) {
	handler := r.GetRoute(c.Method, c.Path)
	if handler != nil {
		handler(c)
	}
}
