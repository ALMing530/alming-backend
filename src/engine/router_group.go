package engine

type RouterGroup struct {
	prifix string
	parent *RouterGroup
	engine *Engine
}

func (rg *RouterGroup) Group(prifix string) *RouterGroup {
	engine := rg.engine
	newRg := &RouterGroup{
		prifix: prifix,
		parent: rg,
		engine: engine,
	}
	engine.gruops = append(engine.gruops, newRg)
	return newRg
}

func (rg *RouterGroup) AddRoute(method string, path string, handler HandleFunc) {
	rg.engine.router.AddRoute(method, rg.prifix+path, handler)
}

func (rg *RouterGroup) GET(path string, handler HandleFunc) {
	rg.AddRoute("GET", path, handler)
}

func (rg *RouterGroup) POST(path string, handler HandleFunc) {
	rg.AddRoute("POST", path, handler)
}

func (rg *RouterGroup) PUT(path string, handler HandleFunc) {
	rg.AddRoute("PUT", path, handler)
}

func (rg *RouterGroup) DELETE(path string, handler HandleFunc) {
	rg.AddRoute("DELETE", path, handler)
}
