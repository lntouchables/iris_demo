package main

import (
	"github.com/kataras/iris/v12"
	"iris_project/test"
)


func InitRoutes(app *iris.Application) {
	rootRouter := app.Party("/api")
	test.RegisterRoutes(rootRouter)
}


