package middleware

import (
	"compress/gzip"
	"github.com/gomi/iType"
	"io"
	"net/http"
	"strings"
)

var gzipDefaultLevel = -1
var gzipSchema = "gzip"

type compressWriter struct {
	http.ResponseWriter
	io.Writer
}

//Compress ...
func Compress(level int) iType.Middle {
	if level == 0 {
		level = gzipDefaultLevel
	}
	return func(ctx *iType.Ctx, next iType.BindMiddle) error {
		req := ctx.Req
		if strings.Contains(req.Header.Get(iType.HeaderAccpetEncoding), gzipSchema) {
			res := ctx.Res
			defer func() {
				if res.Size == 0 {
					//参考echo issue
					res.Header().Del(iType.HeaderContentEncoding)
					res.Header().Del(iType.HeaderContentLength)
				}
			}()
			w, err := gzip.NewWriterLevel(res, level)
			if err != nil {
				return err
			}
			ctx.Res.Writer = &compressWriter{
				Writer:         w,
				ResponseWriter: res,
			}
		}
		next(ctx)
		return nil
	}
}

//Write ...
func (c *compressWriter) Write(b []byte) (int, error) {
	if c.Header().Get(iType.HeaderContentType) == "" {
		c.Header().Set(iType.HeaderContentType, http.DetectContentType(b))
	}
	return c.ResponseWriter.Write(b)
}

//WriteHeader ...
func (c *compressWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		c.Header().Del(iType.HeaderContentEncoding)
	}

	c.Header().Del(iType.HeaderContentLength)
	c.ResponseWriter.WriteHeader(code)
}

//Flush ...
func (c *compressWriter) Flush() {
	c.Writer.(*gzip.Writer).Flush()
}