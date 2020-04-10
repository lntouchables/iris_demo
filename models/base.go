package models

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"iris/config"
	"iris/sysinit"
	"iris/validates"
	"strconv"
)

func CreateSystemData(perms []*validates.PermissionRequest)  {
	permIds := CreateSystemAdminPermission(perms) //初始化权限
	role := CreateSystemAdminRole(permIds)
	if role.ID != 0{
		CreateSystemAdmin(role.ID)
	}//初始化角色
}

func CreateSystemAdminRole(permIds []uint) *Role {
	rr := &validates.RoleRequest{
		Name:        "admin",
		DisplayName: "管理员",
		Description: "管理员",
	}
	role := NewRoleByStruct(rr)
	role.GetRoleByName()
	if role.ID == 0 {
		role.CreateRole(permIds)
	}

	return role
}

/**
 * 创建系统权限
 * @return
 */
func CreateSystemAdminPermission(perms []*validates.PermissionRequest) []uint {
	var permIds []uint
	for _, perm := range perms {
		p := NewPermission(0, perm.Name, perm.Act)
		p.DisplayName = perm.DisplayName
		p.Description = perm.Description
		p.GetPermissionByNameAct()
		if p.ID != 0 {
			continue
		}
		p.CreatePermission()
		permIds = append(permIds, p.ID)
	}
	return permIds
}


func CreateSystemAdmin(roleId uint)  {
	aul := &validates.CreateUpdateUserRequest{
		Username: config.Config.Admin.UserName,
		Password: config.Config.Admin.Pwd,
		Name:     config.Config.Admin.Name,
		RoleIds:  []uint{roleId},
	}

	user := NewUserByStruct(aul)
	user.GetUserByUsername()
	if user.ID == 0{
		user.CreateUser(aul)
	}
}

func IsNotFound(err error)  {
	if ok := errors.Is(err, gorm.ErrRecordNotFound); !ok && err != nil{
		color.Red(fmt.Sprintf("error :%v \n", err))
	}
}

func GetAll(string, orderBy string, offset, limit int) *gorm.DB {
	db := sysinit.Db
	if len(orderBy) > 0{
		db.Order(orderBy + "desc")
	}else{
		db.Order("create_at desc")
	}

	if len(string) > 0{
		db.Where("name LIKE ?", "%"+string+"%")
	}

	if offset > 0{
		db.Offset((offset - 1) * limit)
	}

	if limit > 0{
		db.Limit(limit)
	}

	return db
}

func DelAllData()  {
	sysinit.Db.Unscoped().Delete(&OauthToken{})
	sysinit.Db.Unscoped().Delete(&Permission{})
	sysinit.Db.Unscoped().Delete(&Role{})
	sysinit.Db.Unscoped().Delete(&User{})
	sysinit.Db.Exec("DELETE FROM casbin_rule;")
}

func Update(v, d interface{}) error {
	if err := sysinit.Db.Model(v).Update(d).Error; err != nil{
		return err
	}
	return nil
}

func GetRolesForUser(uid uint) []string {
	uids, err := sysinit.Enforcer.GetRolesForUser(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		color.Red(fmt.Sprintf("GetRolesForUser 错误: %v", err))
		return []string{}
	}

	return uids
}


func GetPermissionsForUser(uid uint) [][]string {
	return sysinit.Enforcer.GetPermissionsForUser(strconv.FormatUint(uint64(uid), 10))
}

func DropTables() {
	sysinit.Db.DropTable("users", "roles", "permissions", "oauth_tokens", "casbin_rule")
}