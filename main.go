package main

import (
	"fmt"
	"iris/config"
	"time"
	"github.com/kataras/iris/v12"
)


func main() {
	f := NewLogFile()
	defer f.Close()

	api := NewApp()
	api.Logger().SetOutput(f) //记录日志

	if err := api.Run(
		iris.Addr(fmt.Sprintf("%s:%d", config.Config.Host, config.Config.Port)),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(time.RFC3339),
	); err != nil {
		fmt.Println(err)
	}
}