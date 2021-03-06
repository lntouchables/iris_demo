package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"iris/sysinit"
	"iris/validates"
	"strconv"
	"time"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"unique; not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
}

func NewRole(id uint, name string) *Role {
	return &Role{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name: name,
	}
}

func NewRoleByStruct(rr *validates.RoleRequest) *Role {
	return &Role{
		Model:       gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        rr.Name,
		DisplayName: rr.DisplayName,
		Description: rr.Description,
	}
}

func (r *Role) GetRoleById() {
	IsNotFound(sysinit.Db.Where("id = ?", r.ID).First(r).Error)
}

func (r *Role) GetRoleByName() {
	IsNotFound(sysinit.Db.Where("name = ?", r.Name).First(r).Error)
}

func (r *Role) DeleteRoleById() {
	if err := sysinit.Db.Delete(r).Error; err != nil{
		color.Red(fmt.Sprintf("DeleteRoleErr:%s \n", err))
	}
}

func GetAllRoles(name, orderBy string, offset, limit int) (roles []*Role) {
	if err := GetAll(name, orderBy, offset, limit).Find(&roles).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllRoleErr:%s \n", err))
	}
	return
}

func (r *Role) CreateRole(permIds []uint) {
	if err := sysinit.Db.Create(r).Error; err != nil{
		color.Red(fmt.Sprintf("CreateRoleErr:%v \n", err))
	}
	addPerms(permIds, r)
	return
}

func addPerms(permIds []uint, role *Role) {
	if len(permIds) > 0 {
		roleId := strconv.FormatUint(uint64(role.ID), 10)
		if _, err := sysinit.Enforcer.DeletePermissionsForUser(roleId); err != nil {
			color.Red(fmt.Sprintf("AppendPermsErr:%s \n", err))
		}
		var perms []Permission
		sysinit.Db.Where("id in (?)", permIds).Find(&perms)
		for _, perm := range perms {
			if _, err := sysinit.Enforcer.AddPolicy(roleId, perm.Name, perm.Act); err != nil {
				color.Red(fmt.Sprintf("AddPolicy:%s \n", err))
			}
		}
	}
}

func (r *Role) UpdateRole(rj *validates.RoleRequest, permIds []uint) {
	if err := Update(r, rj); err != nil{
		color.Red(fmt.Sprintf("UpdateRoleErr:%s \n", err))
	}

	addPerms(permIds, r)
	return
}


// 角色权限
func (r *Role) RolePermisions() []*Permission {
	perms := GetPermissionsForUser(r.ID)
	var ps []*Permission
	for _, perm := range perms {
		if len(perm) >= 3 && len(perm[1]) > 0 && len(perm[2]) > 0 {
			p := NewPermission(0, perm[1], perm[2])
			p.GetPermissionByNameAct()
			if p.ID > 0 {
				ps = append(ps, p)
			}
		}
	}
	return ps
}