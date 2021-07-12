package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
)

func InitLogger(app *iris.Application) {
	customLogger := logger.New(logger.Config{
		//状态显示状态代码
		Status:true,
		// IP显示请求的远程地址
		IP:true,
		//方法显示http方法
		Method:true,
		// Path显示请求路径
		Path:true,
		// Query将url查询附加到Path。
		Query:true,
		//Columns：true，
		// 如果不为空然后它的内容来自`ctx.Values(),Get("logger_message")
		//将添加到日志中。
		MessageContextKeys:[]string{"logger_message"},
		//如果不为空然后它的内容来自`ctx.GetHeader（“User-Agent”）
		MessageHeaderKeys:[]string{"User-Agent"},
	})
	app.Use(customLogger)
}
