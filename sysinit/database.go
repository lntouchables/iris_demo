package sysinit

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"iris/config"
)

var Db *gorm.DB

func init()  {
	var err error
	var conn string
	if config.Config.DB.Adapter == "postgres" {
		conn = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", config.Config.DB.User, config.Config.DB.Password, config.Config.DB.Host, config.Config.DB.Name)
	}else{
		panic(errors.New("not supported database adapter"))
	}

	Db, err = gorm.Open(config.Config.DB.Adapter, conn)

	if err != nil{
		panic(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return "iris_" + defaultTableName
	}

	Db.DB().SetConnMaxLifetime(10)
	Db.DB().SetMaxOpenConns(100)
}