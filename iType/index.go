package iType

import (
	"net/http"
)

//Ctx ...
type Ctx struct {
	Res *http.ResponseWriter
	Req *http.Request
}

//Middle ...
type Middle func(Ctx, BindMiddle)

//BindMiddle ...
type BindMiddle func(Ctx)

//ExtendMiddleSlice ...
type ExtendMiddleSlice []Middle

//CombineMiddle ...
func CombineMiddle(m ExtendMiddleSlice) BindMiddle {
	return m.Combine(func(next BindMiddle, previous Middle) BindMiddle {
		return func(ctx Ctx) {
			previous(ctx, next)
		}
	})
}

//Combine ...
func (e ExtendMiddleSlice) Combine(callback func(BindMiddle, Middle) BindMiddle) BindMiddle {
	length := len(e)
	if length == 0 {
		return func(ctx Ctx) {
		}
	}

	last := func(ctx Ctx) {
		e[length-1](ctx, func(ctx Ctx) {})
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
