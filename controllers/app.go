package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"iris/libs"
	"iris/models"
	"iris/validates"
	"net/http"
)

func UserLogin(ctx iris.Context)  {
	aul := new(validates.LoginRequest)
	if err := ctx.ReadJSON(aul); err != nil{
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*aul)
	if err != nil{
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans){
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, nil ,e))
				return
			}
		}
	}
	ctx.Application().Logger().Infof("%s 登录系统", aul.Username)
	ctx.StatusCode(iris.StatusOK)
	user := models.NewUser(0, aul.Username)
	user.GetUserByUsername()

	response, status, msg := user.CheckLogin(aul.Password)
	_, _ = ctx.JSON(ApiResource(status, response, msg))
	return
}

func UserLogout(ctx iris.Context)  {
	aui := ctx.Values().GetString("auth_user_id")
	uid := uint(libs.ParseInt(aui, 0))
	models.UserAdminLogout(uid)

	ctx.Application().Logger().Infof("%d 退出系统", uid)
	ctx.StatusCode(http.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "退出"))
}