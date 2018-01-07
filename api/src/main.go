package main

import (
	"controller"

	"github.com/valyala/fasthttp"
)

func main() {
	m := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/user/login":
			controller.UserLogin(ctx)
		case "/user/register":
			controller.UserRegister(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	fasthttp.ListenAndServe(":8081", m)
}
