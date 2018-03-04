package main

import (
	"controller"

	"github.com/valyala/fasthttp"
)

func main() {
	m := func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if err := recover(); err != nil {
				controller.ErrorRes(ctx, err)
			}
		}()
		switch string(ctx.Path()) {
		case "/user/login":
			if string(ctx.Method()) == "POST" {
				controller.UserLogin(ctx)
			} else {
				ctx.Error("method is not allowed", fasthttp.StatusMethodNotAllowed)
			}
		case "/user/register":
			controller.UserRegister(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	fasthttp.ListenAndServe(":8081", m)
}
