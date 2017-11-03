package gomi

import (
	"gomi/iType"
	"net/http"
)

var overall iType.ExtendMiddleSlice

func init() {
}

//App ...
type App struct {
}

//Use ...
func (a *App) Use(m ...iType.Middle) {
	overall = append(overall, m...)
}

//Run
func (a *App) Run(addr string) {
	var passage iType.BindMiddle
	if len(overall) != 0 {
		passage = iType.CombineMiddle(overall)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := iType.Ctx{
			Req: r,
			Res: &w,
		}
		passarge(ctx)
	})
	http.ListenAndServe(addr, nil)
}
