# gomi

## 1.Install 
```go
go get -u github.com/UserNameSoMany/gomi
```

## 2.2. Getting Started

```go

import(
	"fmt"
	"github.com/gomi"
	"github.com/gomi/iType"
	"github.com/gomi/route"
	"github.com/gomi/middleware"
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
		ctx.Res.Write([]byte("hello"))
		return nil
	})
	app.Use(router.Route())
	app.Run(":7890")
```

## 3.Features
* 支持全局中间件
* 支持简单路由组
* 支持路由组中间件
* 支持静态路由和参数路由
* 支持路由多处理函数

## 4.Usage
#### 1) 常规静态路由
路由属于各路由组，路由组可以设置前缀，如果路由组不需要前缀，则传入空字符串
```go
import(
	"fmt"
	"github.com/gomi"
	"github.com/gomi/iType"
	"github.com/gomi/route"
)

func main() {
	app := gomi.New()
	//初始化路由组
	router := route.New("/api/v2")
	router.Get("/a", func(ctx *iType.Ctx, next iType.BindMiddle) error {
		ctx.Res.Write([]byte("hello"))
		return nil
	})
	app.Use(router.Route())
	app.Run(":7890")
```

#### 2) 参数化路由
```go
import(
	"fmt"
	"github.com/gomi"
	"github.com/gomi/iType"
	"github.com/gomi/route"
)
func main() {
	app := gomi.New()
	app.Use(func(ctx *iType.Ctx, next iType.BindMiddle) error {
		fmt.Println("hello, this is global middle")
		return next(ctx)
	})
	//初始化路由组
	router := route.New("/api/v2")
	router.Get("/a/:id", func(ctx *iType.Ctx, next iType.BindMiddle) error {
		//获许路由参数
		pathParam := ctx.GetPathStringParam("id")
		ctx.Res.Write([]byte(pathParam))
		return nil
	})
	app.Use(router.Route())
	app.Run(":7890")
```

#### 3) 请求
支持POST， PUT， DELETE， GET请求

```go
	router.Post("/b",func(ctx *iType.Ctx, next iType.BindMiddle)error {
		//ctx.Req 为http.Request
		fmt.Println(ctx.Req.Header.Get("Content-Type"))
		//获取无文件表单 a
		fmt.Println(ctx.Input.FormValue("a"))

		//获取querystring
		fmt.Println(ctx.Input.QueryStringValue("c"))
		testJSON := struct{
			a string `json:"a"`
		}{}
		
		//ctx.Input.RequestBody为content-type为json的请求的请求体（需使用middle.Parse中间件）
		json.Unmarshal(ctx.Input.RequestBody, &testJSON)
		fmt.Printf("json %v", testJSON)
		
		//ctx.Res为http.ResponseWriter
		ctx.Res.Write([]byte("hellopost"))
		return nil
	})
```
#### 4） 中间件

支持三种级别的中间件：1.全局中间件。 2.路由组中间件 3.路由中间件，即路由支持多个处理函数
中间件和处理函数是同一类型，第一个参数ctx代表着传递的上下文，包含着 http.request 和 http.ResponseWriter
而第二个参数next代表这下一个处理函数，如果希望将控制权交给下一个函数，则调用next(ctx)
基本流程为 全局中间件 -> 路由组中间件 -> 路由处理函数 -> 全局中间件
中间件及其灵活，根据所处位置不同，触发时机不通，基本可以根据Use顺序进行

```go
import(
	"fmt"
	"github.com/gomi"
	"github.com/gomi/iType"
	"github.com/gomi/route"
)
func main() {
	app := gomi.New()
	//使用全局中间件，任何请求，无论是否匹配到路由，都会出发全局中间件
	app.Use(func(ctx *iType.Ctx, next iType.BindMiddle) error {
		fmt.Println("hello, this is global middle")
		return next(ctx)
	})
	//初始化路由组
	router := route.New("/api/v2")
	//路由组中间件，只有当请求路径和路由组中的某个路径匹配上才会执行
	router.Use(func(ctx *iType.Ctx, next iType.BindMiddle) error {
		if next != nil {
			return next(ctx)
		}
		return nil
	})
	//路由支持多个处理
	router.Get("/a/:id", func(ctx *iType.Ctx, next iType.BindMiddle) error {
		//逻辑处理
		
		//将控制权转交给下个处理函数
		return next(ctx)
	}, func(ctx *iType.Ctx, next iType.BindMiddle) error {
		//获许路由参数
		pathParam := ctx.GetPathStringParam("id")
		ctx.Res.Write([]byte(pathParam))

		//如果调用next(ctx)，则控制权会继续向后转交
		return nil
	})

	//将路由组挂载到应用上
	app.Use(router.Route())
	
	//当前一个路由组没有匹配到或者转交了控制权后，触发此处中间件
	app.Use(func(ctx *iType.Ctx, next iType.BindMiddle) error {
		if next != nil {
			//将控制权转交给 /api/v3路由组
			return next(ctx)
		}
		return nil
	})
	//初始化另一路由组
	router2 := route.New("/api/v3")
	//路由支持多个处理，匹配 ／api/v2/a/a
	router.Get("/a/:id", func(ctx *iType.Ctx, next iType.BindMiddle) error {
		//逻辑处理
		
		//将控制权转交给下个处理函数
		return next(ctx)
	}, func(ctx *iType.Ctx, next iType.BindMiddle) error {
		//获许路由参数
		pathParam := ctx.GetPathStringParam("id")
		ctx.Res.Write([]byte(pathParam))
		return nil
	})

	//将路由组挂载到应用上
	app.Use(router2.Route())
 
	//最后执行的中间件，当没有任何路由匹配上，或者路由处理函数转交了控制权，则触发
	//可以处理404
	app.Use(func(ctx *iType.Ctx, next iType.BindMiddle) error {
		if next != nil {
			return next(ctx)
		}
		return nil
	})
	app.Run(":7890")
```








	router.Get("/a/:id", func(ctx *iType.Ctx, next iType.BindMiddle) error {
    		ctx.Res.Write([]byte(ctx.GetPathStringParam("id")))
    		return nil
    })
    router.Get("/a/:id/c/:t", func(ctx *iType.Ctx, next iType.BindMiddle) error {
    		ctx.Res.Write([]byte(ctx.GetPathStringParam("id")))
    		ctx.Res.Write([]byte(ctx.GetPathStringParam("t")))
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

}
/*

支持get put delete post
不支持参数化url

only supprot get,put,delete and post method
not support url params
*/
