package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	gf "github.com/snowlyg/gotransformer"
	"iris/models"
	"iris/transformer"
	"iris/validates"
	"strconv"
	"time"
)

func GetProfile(ctx iris.Context) {
	userId := ctx.Values().Get("auth_user_id").(uint)
	user := models.NewUser(userId, "")
	user.GetUserById()
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, userTransform(user), ""))
}

func GetUser(ctx iris.Context)  {
	id, _ := ctx.Params().GetUint("id")
	user := models.NewUser(id, "")
	user.GetUserById()
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, userTransform(user), "操作成功"))
}

func usersTransform(users []*models.User) []*transformer.User {
	var us []*transformer.User
	for _, user := range users {
		u := userTransform(user)
		us = append(us, u)
	}
	return us
}

func CreateUser(ctx iris.Context)  {
	aul := new(validates.CreateUpdateUserRequest)
	if err := ctx.ReadJSON(aul); err != nil{
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*aul)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, nil, e))
				return
			}
		}
	}

	user := models.NewUserByStruct(aul)
	user.CreateUser(aul)
	ctx.StatusCode(iris.StatusOK)
	if user.ID == 0 {
		_, _ = ctx.JSON(ApiResource(false, user, "操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ApiResource(true, nil, "操作成功"))
		return
	}
}

func UpdateUser(ctx iris.Context) {
	aul := new(validates.CreateUpdateUserRequest)

	if err := ctx.ReadJSON(aul); err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
	}

	err := validates.Validate.Struct(*aul)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans) {
			if len(e) > 0 {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	user := models.NewUser(id, "")
	if user.Username == "username" {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, "不能编辑管理员"))
		return
	}

	user.UpdateUser(aul)
	ctx.StatusCode(iris.StatusOK)
	if user.ID == 0 {
		_, _ = ctx.JSON(ApiResource(false, user, "操作失败"))
		return
	} else {
		_, _ = ctx.JSON(ApiResource(true, nil, "操作成功"))
		return
	}
}


func DeleteUser(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")

	user := models.NewUser(id, "")
	if user.Username == "username" {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, "不能删除管理员"))
		return
	}

	user.DeleteUser()

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "删除成功"))
}

func GetAllUsers(ctx iris.Context) {
	offset := ctx.URLParamIntDefault("offset", 1)
	limit := ctx.URLParamIntDefault("limit", 15)
	name := ctx.URLParam("name")
	orderBy := ctx.URLParam("orderBy")

	users := models.GetAllUsers(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, usersTransform(users), "操作成功"))
}


func userTransform(user *models.User) *transformer.User {
	u := &transformer.User{}
	g := gf.NewTransform(u, user, time.RFC3339)
	_ = g.Transformer()

	roleIds := models.GetRolesForUser(user.ID)
	var ris []int
	var roleName string
	for num, roleId := range roleIds {
		ri, _ := strconv.Atoi(roleId)
		ris = append(ris, ri)
		role := models.NewRole(uint(ri), "")
		role.GetRoleById()
		if num == len(roleIds)-1 {
			roleName += role.DisplayName
		} else {
			roleName += role.DisplayName + ","
		}
	}
	u.RoleIds = ris
	u.RoleName = roleName
	return u
}
