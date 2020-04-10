package models

import (
	"github.com/jinzhu/gorm"
	"iris/sysinit"
)

type OauthToken struct {
	gorm.Model
	Token     string `gorm:"not null default '' comment('Token') VARCHAR(191)"`
	UserId    uint   `gorm:"not null default '' comment('UserId') VARCHAR(191)"`
	Secret    string `gorm:"not null default '' comment('Secret') VARCHAR(191)"`
	ExpressIn int64  `gorm:"not null default 0 comment('是否是标准库') BIGINT(20)"`
	Revoked   bool
}

type Token struct {
	Token string `json:"access_token"`
}

func (ot *OauthToken) OauthTokenCreate() *Token {
	sysinit.Db.Create(ot)
	return &Token{ot.Token}
}

func (ot *OauthToken) GetOauthTokenByToken(token string) {
	sysinit.Db.Where("token = ?", token).First(&ot)
}

func (ot *OauthToken) UpdateOauthTokenByUserId(userId uint) {
	sysinit.Db.Model(ot).Where("revoked = ?", false).Where("user_id = ?", userId).
		Update(map[string]interface{}{"revoked": true})
}
