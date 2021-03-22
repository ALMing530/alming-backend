package engine

import "strings"

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
func (r *router) GetRoute(method string, path string) (node *trieNode, params map[string]string) {
	route, ok := r.routes[method]
	if !ok {
		return nil, nil
	}
	node = route.search(path)
	if node != nil {
		params = make(map[string]string)
		patternSplit := patternToParts(node.pattern)
		pathSplit := patternToParts(path)
		for index, item := range patternSplit {
			if item[0] == ':' {
				params[item[1:]] = pathSplit[index]
			} else if item[0] == '*' {
				params["*"] = strings.Join(pathSplit[index:], "/")
			}
		}
		return node, params
	}
	return nil, nil
}
func (r *router) handle(c *Context) {
	routeNode, params := r.GetRoute(c.Method, c.Path)
	c.PathParams = params
	if routeNode != nil {
		key := c.Method + ":" + routeNode.pattern
		r.handler[key](c)
	}
}
