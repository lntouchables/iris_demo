package controllers

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	gf "github.com/snowlyg/gotransformer"
	"iris/libs"
	"iris/models"
	"iris/transformer"
	"iris/validates"
	"time"
)

func GetRole(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	role := models.NewRole(id, "")
	role.GetRoleById()
	ctx.StatusCode(iris.StatusOK)
	rr := roleTransform(role)
	rr.Perms = permsTransform(role.RolePermisions())
	_, _ = ctx.JSON(ApiResource(true, rr, "操作成功"))
}

func CreateRole(ctx iris.Context)  {
	roleJson := new(validates.RoleRequest)
	if err := ctx.ReadJSON(roleJson); err != nil{
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*roleJson)
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

	role := models.NewRoleByStruct(roleJson)
	role.CreateRole(roleJson.PermissionsIds)
	ctx.StatusCode(iris.StatusOK)
	if role.ID == 0{
		_, _ = ctx.JSON(ApiResource(false, role, "操作失败"))
	}else{
		_, _ = ctx.JSON(ApiResource(true, nil, "操作成功"))
		return
	}
}

func UpdateRole(ctx iris.Context)  {
	roleForm := new(validates.RoleRequest)
	if err := ctx.ReadJSON(roleForm); err != nil{
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, err.Error()))
		return
	}

	err := validates.Validate.Struct(*roleForm)
	if err != nil{
		errs := err.(validator.ValidationErrors)
		for _, e := range errs.Translate(validates.ValidateTrans){
			if len(e) > 0{
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, nil, e))
				return
			}
		}
	}

	id, _ := ctx.Params().GetUint("id")
	role := models.NewRole(id, "")
	role.GetRoleById()
	if role.Name == "admin" {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, nil, "不能编辑管理员角色"))
		return
	}

	roleJson := new(validates.RoleRequest)
	roleJson.Name = roleForm.Name
	roleJson.Description = roleForm.Description
	roleJson.DisplayName = roleForm.DisplayName

	role.UpdateRole(roleJson, roleForm.PermissionsIds)
	ctx.StatusCode(iris.StatusOK)
	if role.ID == 0{
		_, _ = ctx.JSON(ApiResource(false, role, "操作失败"))
		return
	}else{
		_, _ = ctx.JSON(ApiResource(true, nil, "操作失败"))
		return
	}
}

func DeleteRole(ctx iris.Context)  {
	id, _ := ctx.Params().GetUint("id")
	role := models.NewRole(id, "")
	role.GetRoleById()
	if role.Name == "admin"{
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(true, nil, "不能删除管理员角色"))
		return
	}
	role.DeleteRoleById()

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "删除成功"))
}

func GetAllRoles(ctx iris.Context)  {
	offset := libs.ParseInt(ctx.FormValue("offer"), 1)
	limit := libs.ParseInt(ctx.FormValue("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	roles := models.GetAllRoles(name, orderBy, offset, limit)
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, rolesTransform(roles), "操作成功"))
}

func rolesTransform(roles []*models.Role) []*transformer.Role {
	var rs []*transformer.Role
	for _, role := range roles {
		r := roleTransform(role)
		rs = append(rs, r)
	}
	return rs
}

func roleTransform(role *models.Role) *transformer.Role {
	r := &transformer.Role{}
	g := gf.NewTransform(r, role, time.RFC3339)
	_ = g.Transformer()
	return r
}
