## 持续更新中
### 路由支持路由组
### 目前支持全局中间件，路由组级别中间件，路由级别中间件
## TODO
### URL增加参数化支持
package main

import(
	"fmt"
	"gomi"
	"gomi/iType"
	"gomi/route"
	"gomi/middleware"
	"encoding/json"
)

func main() {
	app := gomi.New()
	app.Use(func(ctx *iType.Ctx, next iType.BindMiddle) error {
		fmt.Println("hello, this is global middle")
		return next(ctx)
	})


	//support prefix,if not use, please use ""
	router := route.New("/api/v2")

	//route middle
	router.Use(middleware.Parse)
	router.Get("/a", func(ctx *iType.Ctx, next iType.BindMiddle) error {
		fmt.Println("hello")
		return next(ctx)
	}, func(ctx *iType.Ctx, next iType.BindMiddle) error {
		ctx.Res.Write([]byte("hello"))
		return nil
	})
	router.Post("/b",func(ctx *iType.Ctx, next iType.BindMiddle)error {
		fmt.Println(ctx.Req.Header.Get("Content-Type"))
		fmt.Println(ctx.Input.FormValue("a"))
		fmt.Println(ctx.Input.QueryStringValue("c"))
		testJSON := struct{
			a string `json:"a"`
		}{}

		json.Unmarshal(ctx.Input.RequestBody, &testJSON)
		fmt.Printf("json %v", testJSON)

		ctx.Res.Write([]byte("hellopost"))
		return nil
	})
	router.Put("/a", func(ctx *iType.Ctx, next iType.BindMiddle)error {
		ctx.Res.Write([]byte("helloput"))
		return nil
	})
	router.Delete("/a", func(ctx *iType.Ctx, next iType.BindMiddle)error {
		ctx.Res.Write([]byte("hellodelete"))
		return nil
	})

	//router2
	router2 := route.New("/api/v3")

	//router2 middle
	router2.Use(middleware.Parse)
	router2.Get("/a", func(ctx *iType.Ctx, next iType.BindMiddle)error {
		ctx.Res.Write([]byte("hellov3"))
		return nil
	})
	app.Use(router.Route())
	app.Use(router2.Route())
	app.Run(":7890")
}
/*

支持get put delete post
不支持参数化url

only supprot get,put,delete and post method
not support url params
*/
