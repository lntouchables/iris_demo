package models

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/iris-contrib/middleware/jwt"
	"github.com/jameskeane/bcrypt"
	"github.com/jinzhu/gorm"
	"iris/libs"
	"iris/sysinit"
	"iris/validates"
	"strconv"
	"time"
)

type User struct {
	gorm.Model
	Name string `gorm:"not null VARCHAR(191)"`
	Username string `gorm:"unique;VARCHAR(191)"`
	Password string `gorm:"not null VARCHAR(191)"`
}

func NewUser(id uint, username string) *User {
	return &User{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username:username,
	}
}

func NewUserByStruct(ru *validates.CreateUpdateUserRequest) *User {
	return &User{
		Model: gorm.Model{
			ID:        0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: ru.Username,
		Name: ru.Name,
		Password: libs.HashPassword(ru.Password),
	}
}

func (u *User) GetUserByUsername() {
	IsNotFound(sysinit.Db.Where("username = ?", u.Username).First(u).Error)
}

func (u *User) GetUserById() {
	IsNotFound(sysinit.Db.Where("id = ?", u.ID).First(u).Error)
}

func (u *User) DeleteUser() {
	if err := sysinit.Db.Delete(u).Error; err != nil{
		color.Red(fmt.Sprintf("DeleteUserByIdErr:%s \n", err))
	}
}

func GetAllUsers(name, orderBy string, offset, limit int) []*User {
	var users []*User
	q := GetAll(name, orderBy, offset, limit)
	if err := q.Find(&users).Error; err != nil{
		color.Red(fmt.Sprintf("GetAllUserErr:%s \n", err))
		return nil
	}
	return users
}

func (u *User) CreateUser(aul *validates.CreateUpdateUserRequest) {
	u.Password = libs.HashPassword(aul.Password)
	if err := sysinit.Db.Create(u).Error; err != nil {
		color.Red(fmt.Sprintf("CreateUserErr:%s \n", err))
	}
	addRoles(aul, u)
	return
}

func (u *User) UpdateUser(uj *validates.CreateUpdateUserRequest) {
	uj.Password = libs.HashPassword(uj.Password)
	if err := Update(u, uj); err != nil{
		color.Red(fmt.Sprintf("UpdateUserErr:%s \n ", err))
	}
	addRoles(uj, u)
}

func addRoles(uj *validates.CreateUpdateUserRequest, user *User)  {
	if len(uj.RoleIds) > 0{
		userId := strconv.FormatUint(uint64(user.ID), 10)
		if _, err := sysinit.Enforcer.DeleteRolesForUser(userId); err != nil{
			color.Red(fmt.Sprintf("CreateUserErr:%s \n", err))
		}

		for _, roleId := range uj.RoleIds{
			roleId := strconv.FormatUint(uint64(roleId), 10)
			if _, err := sysinit.Enforcer.AddRoleForUser(userId, roleId); err != nil{
				color.Red(fmt.Sprintf("CreateUserErr:%s \n", err))
			}
		}
	}
}

func (u *User) CheckLogin(password string) (*Token, bool, string) {
	if u.ID == 0{
		return nil, false, "用户不存在"
	}else{
		if ok := bcrypt.Match(password, u.Password); ok{
			token := jwt.NewTokenWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
				"exp" : time.Now().Add(time.Hour * time.Duration(1)).Unix(),
				"iat": time.Now().Unix(),
			})
			tokenString, _ := token.SignedString([]byte("HS2JDFKhu7Y1av7b"))
			oauthToken := new(OauthToken)
			oauthToken.Token = tokenString
			oauthToken.UserId = u.ID
			oauthToken.Secret = "secret"
			oauthToken.Revoked = false
			oauthToken.ExpressIn = time.Now().Add(time.Hour * time.Duration(1)).Unix()
			oauthToken.CreatedAt = time.Now()

			response := oauthToken.OauthTokenCreate()
			return response, true, "登录成功"
		}else{
			return nil, false, "用户名或密码错误"
		}
	}
}

func UserAdminLogout(userId uint) bool {
	ot := OauthToken{}
	ot.UpdateOauthTokenByUserId(userId)
	return ot.Revoked
}