package main

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"time"
)

func main() {
	app := iris.New()
	InitLogger(app)
	InitRoutes(app)
	if err := app.Run(
		iris.Addr(":8089"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(time.RFC3339),
	); err != nil {
		fmt.Println(err)
	}
}
