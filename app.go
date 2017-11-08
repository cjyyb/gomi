package gomi

import (
	"gomi/iType"
	"net"
	"net/http"

	"github.com/golang/glog"
)

var overall iType.ExtendMiddleSlice

func init() {
}

//App ...
type App struct {
	server *http.Server
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var passage iType.BindMiddle
	if len(overall) != 0 {
		passage = iType.CombineMiddle(overall)
	}
	ctx := iType.New(r, w)
	defer func() {
		//do your self err handle
		err := recover()
		if err != nil {
			glog.Errorf("main process catch error, err: %v", err)
		}
	}()
	defer iType.Release(ctx)
	passage(ctx)
}

//New ...
func New() *App {
	app := App{
		server: new(http.Server),
	}
	app.server.Handler = &app
	return &app
}

//Use ...
func (a *App) Use(m ...iType.Middle) {
	overall = append(overall, m...)
}

//Run ...
func (a *App) Run(addr string) error {
	l, err := newListener(addr)
	if err != nil {
		glog.Errorln("Start http server failed, err: %v", err)
		panic(err)
	}
	return a.server.Serve(l)
}

func newListener(addr string) (net.Listener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return l, nil
}
