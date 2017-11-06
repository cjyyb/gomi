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
	query url.Values
}

//QueryString ...
func (c *Ctx) QueryString(name string) string {
	return c.Req.URL.RawQuery
}

//QueryStringValue ...
func (c *Ctx) QueryStringValue(name string) string {
	if c.query == nil {
		c.query = c.Req.URL.Query()
	}
	return c.query.Get(name)
}

//QueryIntValue ...
func (c *Ctx) QueryIntValue(name string) (int, error) {
	if c.query == nil {
		c.query = c.Req.URL.Query()
	}
	value := c.query.Get(name)
	if value == "" {
		return 0, nil
	}
	iv, err := strconv.Atoi(value)
	return iv, err
}

//CtxPool ...
var CtxPool = sync.Pool{
	New: func() interface{} {
		return new(Ctx)
	},
}

//Middle ...
type Middle func(*Ctx, BindMiddle)

//BindMiddle ...
type BindMiddle func(*Ctx)

//ExtendMiddleSlice ...
type ExtendMiddleSlice []Middle

//CombineMiddle ...
func CombineMiddle(m ExtendMiddleSlice) BindMiddle {
	return m.combine(func(next BindMiddle, previous Middle) BindMiddle {
		return func(ctx *Ctx) {
			previous(ctx, next)
		}
	})
}

//Combine ...
func (e ExtendMiddleSlice) combine(callback func(BindMiddle, Middle) BindMiddle) BindMiddle {
	length := len(e)
	if length == 0 {
		return func(ctx *Ctx) {
		}
	}

	last := func(ctx *Ctx) {
		e[length-1](ctx, func(ctx *Ctx) {})
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
