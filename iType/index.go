package iType

import (
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

//Ctx ...
type Ctx struct {
	Res   http.ResponseWriter
	Req   *http.Request
	URL   URL
	Input *Input
}

//GetStringParams ...
func (c *Ctx) GetPathStringParams(name string) string {
	return c.URL.GetPathStringParams(name)
}

//URL ...
type URL struct {
	Params map[string]string
}

//GetStringParams ...
func (u *URL) GetPathStringParams(name string) string {
	return u.Params[name]
}

//Input ...
type Input struct {
	RequestBody []byte
	Ctx         *Ctx
	query       url.Values
}

//QueryString ...
func (i *Input) QueryString() string {
	return i.Ctx.Req.URL.RawQuery
}

//QueryStringValue ...
func (i *Input) QueryStringValue(name string) string {
	if i.query == nil {
		i.query = i.Ctx.Req.URL.Query()
	}
	return i.query.Get(name)
}

//FormValue ...
func (i *Input) FormValue(name string) string {
	return i.Ctx.Req.FormValue(name)
}

//QueryIntValue ...
func (i *Input) QueryIntValue(name string) (int, error) {
	if i.query == nil {
		i.query = i.Ctx.Req.URL.Query()
	}
	value := i.query.Get(name)
	if value == "" {
		return 0, nil
	}
	iv, err := strconv.Atoi(value)
	return iv, err
}

//CtxPool ...
var ctxPool = sync.Pool{
	New: func() interface{} {
		return new(Ctx)
	},
}

//New ...
func New(req *http.Request, res http.ResponseWriter) *Ctx {
	ctx := ctxPool.Get().(*Ctx)
	ctx.Res = res
	ctx.Req = req
	ctx.Input = new(Input)
	ctx.Input.Ctx = ctx
	return ctx
}

//Release ...
func Release(ctx *Ctx) {
	ctxPool.Put(ctx)
}

//Middle ...
type Middle func(*Ctx, BindMiddle) error

//BindMiddle ...
type BindMiddle func(*Ctx) error

//ExtendMiddleSlice ...
type ExtendMiddleSlice []Middle

//CombineMiddle ...
func CombineMiddle(m ExtendMiddleSlice) BindMiddle {
	return m.combine(func(next BindMiddle, previous Middle) BindMiddle {
		return func(ctx *Ctx) error {
			return previous(ctx, next)
		}
	})
}

//Combine ...
func (e ExtendMiddleSlice) combine(callback func(BindMiddle, Middle) BindMiddle) BindMiddle {
	length := len(e)
	if length == 0 {
		return func(ctx *Ctx) error {
			return nil
		}
	}

	last := func(ctx *Ctx) error {
		return e[length-1](ctx, func(ctx *Ctx) error {
			return nil
		})
	}

	if length == 1 {
		return last
	}
	m := last
	for i := length - 2; i >= 0; i-- {
		m = callback(m, e[i])
	}
	return m
}
