package test

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)


func demo(ctx iris.Context)  {
	ctx.JSON(iris.Map{
		"message": "demo",
	})
}

func RegisterRoutes(rootRouter router.Party) {
	nodeRouter := rootRouter.Party("/test")

	nodeRouter.Get("/demo", demo)
}
