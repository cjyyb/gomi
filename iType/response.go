package iType

import (
	"net/http"
)

//Response ...
type Response struct {
	Ctx      *Ctx
	Writer   http.ResponseWriter
	Size     int
	IsCommit bool
}

//NewResponse ...
func NewResponse(ctx *Ctx, rw http.ResponseWriter) *Response {
	return &Response{
		Ctx:      ctx,
		Writer:   rw,
		Size:     0,
		IsCommit: false,
	}
}

//Write
func (r *Response) Write(b []byte) (n int, err error) {
	if !r.IsCommit {
		r.WriteHeader(http.StatusOK)
	}
	n, err = r.Writer.Write(b)
	r.Size += n
	return
}

//Header ...
func (r *Response) Header() http.Header {
	return r.Writer.Header()
}

//WriteHeader ...
func (r *Response) WriteHeader(code int) {
	if r.IsCommit {
		return
	}
	r.Writer.WriteHeader(code)
	r.IsCommit = true
}
