package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	"iris/sysinit"
	"iris/validates"
	"time"
)

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null VARCHAR(191)"`
	DisplayName string `gorm:"VARCHAR(191)"`
	Description string `gorm:"VARCHAR(191)"`
	Act         string `gorm:"VARCHAR(191)"`
}

func NewPermission(id uint, name, act string) *Permission {
	return &Permission{
		Model:       gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        name,
		Act:         act,
	}
}

func NewPermissionByStruct(jp *validates.PermissionRequest) *Permission {
	return &Permission{
		Model:       gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:        jp.Name,
		DisplayName: jp.DisplayName,
		Description: jp.Description,
		Act:         jp.Act,
	}
}

func (p *Permission) GetPermissionById() {
	IsNotFound(sysinit.Db.Where("id = ?", p.ID).First(p).Error)
}

func (p *Permission) GetPermissionByNameAct() {
	IsNotFound(sysinit.Db.Where("name = ?", p.Name).Where("act = ?", p.Act).First(p).Error)
}


func (p *Permission) DeletePermissionById() {
	if err := sysinit.Db.Delete(p).Error; err != nil {
		color.Red(fmt.Sprintf("DeletePermissionByIdError:%s \n", err))
	}
}


func GetAllPermissions(name, orderBy string, offset, limit int) (permissions []*Permission) {
	if err := GetAll(name, orderBy, offset, limit).Find(&permissions).Error; err != nil {
		color.Red(fmt.Sprintf("GetAllPermissionsError:%s \n", err))
	}

	return
}


func (p *Permission) CreatePermission() {
	if err := sysinit.Db.Create(p).Error; err != nil {
		color.Red(fmt.Sprintf("CreatePermissionError:%s \n", err))
	}
	return
}


func (p *Permission) UpdatePermission(pj *validates.PermissionRequest) {
	if err := Update(p, pj); err != nil {
		color.Red(fmt.Sprintf("UpdatePermissionError:%s \n", err))
	}
}