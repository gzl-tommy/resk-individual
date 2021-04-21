package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	app := iris.New()
	app.Get("/hello", func(ctx context.Context) {
		ctx.WriteString("hello,world! iris")
	})

	v1 := app.Party("/v1")
	v1.Use(func(c context.Context) {
		logrus.Info("自定义中间件")
		c.Next()
	})

	v1.Get("/user/{id:uint64 main(2)}", func(c context.Context) {
		id := c.Params().GetUint64Default("id", 0)
		c.WriteString(fmt.Sprintf("%d", id))
	})

	v1.Get("/orders/{action:string prefix(a_)}", func(c context.Context) {
		a := c.Params().Get("action")
		c.WriteString(a)
	})

	app.OnAnyErrorCode(func(c context.Context) {
		c.WriteString("看起来服务器出错了耶！")
	})

	app.OnErrorCode(http.StatusNotFound, func(c context.Context) {
		c.WriteString("访问路径不存在哦！")
	})

	err := app.Run(iris.Addr(":8082"))
	if err != nil {
		fmt.Println(err)
	}
}
