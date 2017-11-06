package middleware

import (
	"io"
	"io/ioutil"

	"gomi/iType"
	"gomi/route"

	"github.com/golang/glog"
)

var (
	//DEFAULTMAXMEMORY ...
	DEFAULTMAXMEMORY int64 = 1 << 26 //64MB
)

//Parse ...
func Parse(ctx *iType.Ctx, next iType.BindMiddle) error {
	req := ctx.Req
	if req.Method == route.POST {
		if isJSON(ctx) == true {
			err := parseRequestBody(ctx.Input, ctx.Req.Body, DEFAULTMAXMEMORY)
			if err != nil {
				glog.Errorln("Parse request body failed, err: %v\n", err)
				return err
			}
		}
	}
	return next(ctx)
}

func parseRequestBody(i *iType.Input, reader io.Reader, maxMemory int64) error {
	r := io.LimitReader(reader, maxMemory)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		glog.Errorf("Read requst body failed, err: %v\n", err)
		return err
	}
	i.RequestBody = data
	if err != nil {
		glog.Errorf("Close Request.Body failed, err: %v\n", err)
		return err
	}
	return nil
}

func isJSON(ctx *iType.Ctx) bool {
	req := ctx.Req
	if req.Header.Get("Content-Type") == "application/json" {
		return true
	}
	return false
}
