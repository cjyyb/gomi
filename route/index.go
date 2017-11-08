package route

import (
	"github.com/gomi/iType"
)

var (
	//POST ...
	POST = "POST"

	//PUT ...
	PUT = "PUT"

	//GET ...
	GET = "GET"

	//DELETE ...
	DELETE = "DELETE"
)

type kind uint8

const (
	ukind uint8 = iota
	pkind
)

//Router ...
type Router struct {
	prefix string
	middle []iType.Middle
	route  *route
}

type route struct {
	label     byte
	children  childrenRoute
	parent    *route
	prefix    string
	handler   *handler
	root      *Router
	urlParams []string
	urlValue  []string
	kind      uint8
}

type handler struct {
	get    iType.BindMiddle
	post   iType.BindMiddle
	put    iType.BindMiddle
	delete iType.BindMiddle
}

type childrenRoute []*route

//New ...
func New(prefix string) *Router {
	router := &Router{
		prefix: prefix,
	}
	route := route{
		root: router,
	}
	router.route = &route
	return router
}

//Use ...
func (r *Router) Use(handler iType.Middle) {
	r.middle = append(r.middle, handler)
}

//Get ...
func (r *Router) Get(path string, handler ...iType.Middle) {
	r.route.process(GET, path, handler...)
}

//Post ...
func (r *Router) Post(path string, handler ...iType.Middle) {
	r.route.process(POST, path, handler...)
}

//Put ...
func (r *Router) Put(path string, handler ...iType.Middle) {
	r.route.process(PUT, path, handler...)
}

//Delete ...
func (r *Router) Delete(path string, handler ...iType.Middle) {
	r.route.process(DELETE, path, handler...)
}

//Route ...
func (r *Router) Route() iType.Middle {
	return func(ctx *iType.Ctx, bind iType.BindMiddle) error {
		handler := r.search(ctx)
		if handler == nil {
			return bind(ctx)
		}
		err := handler(ctx)
		if err != nil {
			return err
		}
		return bind(ctx)
	}
}

func (r *Router) search(ctx *iType.Ctx) iType.BindMiddle {
	req := ctx.Req
	path := req.URL.Path
	prefix := r.prefix
	preLength := len(prefix)
	if prefix == "" {
		rr := findHandlerByMethodAndPath(r.route, req.Method, path)
		if rr == nil {
			return nil
		}
		ctx.URL.Params = params(rr.urlParams, rr.urlValue)
		return rr.getHandlerByMethod(req.Method)
	}
	l := 0
	for ; l < preLength && prefix[l] == path[l]; l++ {
	}
	if l != preLength {
		return nil
	}
	rr := findHandlerByMethodAndPath(r.route, req.Method, path[l:])
	if rr == nil {
		return nil
	}
	ctx.URL.Params = params(rr.urlParams, rr.urlValue)
	return rr.getHandlerByMethod(req.Method)
}

func params(urlParams, urlValue []string) map[string]string {
	ptov := map[string]string{}
	if len(urlParams) > 0 {
		vLength := len(urlValue)
		for i, value := range urlParams {
			if i < vLength {
				ptov[value] = urlValue[i]
			}
		}
	}
	return ptov
}

func findHandlerByMethodAndPath(r *route, method, path string) *route {
	var (
		pvalue        = []string{}
		previousRoute *route
	)
	for {
		if r == nil {
			return nil
		}
		if path == "" {
			goto End
		}
		l := 0
		preLength := 0
		prefix := r.prefix
		if r.label != ':' {
			preLength = len(prefix)
			pathLength := len(path)
			max := preLength
			if max > pathLength {
				max = pathLength
			}
			for ; l < max && prefix[l] == path[l]; l++ {
			}
		}
		if l == preLength {
			path = path[l:]
		} else {
			r = previousRoute
			goto Params
		}
		if path == "" {
			goto End
		}
		if cr := findRouteByLabelAndKind(r, path[0], ukind); cr != nil {
			previousRoute = r
			r = cr
			continue
		}

	Params:
		if cr := findRouteByKind(r, pkind); cr != nil {
			i, l := 0, len(path)
			for ; i < l && path[i] != '/'; i++ {
			}
			pvalue = append(pvalue, path[:i])
			path = path[i:]
			r = cr
			continue
		}
		return nil

	}
End:
	r.urlValue = pvalue
	return r
}

func (r *route) process(method, path string, handler ...iType.Middle) {
	if path[0] != '/' {
		path = "/" + path
	}
	if len(handler) == 0 {
		handler = append(handler, func(ctx *iType.Ctx, next iType.BindMiddle) error {
			return nil
		})
	}
	pname := []string{}
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1
			r.add(method, path[:i], ukind, nil)
			for ; i < l && path[i] != '/'; i++ {
			}
			pname = append(pname, path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)
			if i == l {
				r.add(method, path, pkind, pname, handler...)
				return
			}
			r.add(method, path[:i], pkind, pname)
		}
	}
	r.add(method, path, ukind, pname, handler...)
}

func (r *route) add(method, path string, kind uint8, pname []string, handler ...iType.Middle) {
	if len(handler) != 0 {
		handler = append(r.root.middle, handler...)
	}
	middle := iType.ExtendMiddleSlice(handler)
	for {
		prefix := r.prefix
		prefixLength := len(prefix)
		pathLength := len(path)
		max := prefixLength
		if max > pathLength {
			max = pathLength
		}
		l := 0
		for l = 0; l < max && prefix[l] == path[l]; l++ {
		}
		if l == 0 {
			r.prefix = path
			r.children = nil
			r.label = path[0]
			r.kind = kind
			r.urlParams = pname
			r.addMethodHandler(method, middle)
		} else if l < prefixLength {
			newPrefix := path[0:l]
			otn := convertToNew(r.prefix[l:], r.urlParams, kind, r.children, r.handler)
			r.prefix = newPrefix
			r.children = nil
			r.handler = nil
			r.addChildren(otn)
			r.kind = kind
			r.urlParams = pname
			path = path[l:]
			if l == pathLength {
				r.addMethodHandler(method, middle)
			} else {
				newRoute := &route{
					label:     path[0],
					prefix:    path,
					kind:      kind,
					urlParams: pname,
				}
				newRoute.addMethodHandler(method, middle)
				r.addChildren(newRoute)
			}
		} else if l < pathLength {
			path = path[l:]
			c := findRouteByLabel(r, path[0])
			if c != nil {
				r = c
				continue
			}
			newRoute := &route{
				label:     path[0],
				prefix:    path,
				kind:      kind,
				urlParams: pname,
			}
			newRoute.addMethodHandler(method, middle)
			r.addChildren(newRoute)
		} else {
			r.addMethodHandler(method, middle)
		}
		return
	}
}

func findRouteByLabelAndKind(r *route, label byte, kind uint8) *route {
	for i, value := range r.children {
		if value.label == label && value.kind == kind {
			return r.children[i]
		}
	}
	return nil
}

func findRouteByLabel(r *route, label byte) *route {
	for i, value := range r.children {
		if value.label == label {
			return r.children[i]
		}
	}
	return nil
}

func findRouteByKind(r *route, kind uint8) *route {
	for _, value := range r.children {
		if value.kind == kind {
			return value
		}
	}

	return nil
}

func convertToNew(prefix string, pname []string, kind uint8, children childrenRoute, h *handler) *route {
	if h == nil {
		h = new(handler)
	}
	router := route{
		label:     prefix[0],
		prefix:    prefix,
		children:  children,
		urlParams: pname,
		kind:      kind,
		handler:   h,
	}
	return &router
}

func (r *route) addChildren(c *route) {
	r.children = append(r.children, c)
}

func (r *route) getHandlerByMethod(method string) iType.BindMiddle {
	switch method {
	case GET:
		return r.handler.get
	case POST:
		return r.handler.post
	case PUT:
		return r.handler.put
	case DELETE:
		return r.handler.delete
	}
	return nil
}

func (r *route) addMethodHandler(method string, m iType.ExtendMiddleSlice) {
	if m == nil || len(m) == 0 {
		return
	}
	bm := iType.CombineMiddle(m)
	if r.handler == nil {
		r.handler = new(handler)
	}
	switch method {
	case POST:
		r.handler.post = bm
	case PUT:
		r.handler.put = bm
	case GET:
		r.handler.get = bm
	case DELETE:
		r.handler.delete = bm
	}
}
