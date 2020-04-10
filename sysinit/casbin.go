package sysinit

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"iris/config"
	"path/filepath"
)

var Enforcer *casbin.Enforcer

func init() {
	var err error
	var conn string
	if config.Config.DB.Adapter == "postgres" {
		conn = fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", config.Config.DB.User, config.Config.DB.Password, config.Config.DB.Host, config.Config.DB.Name)
	} else {
		panic(errors.New("not supported database adapter"))
	}

	if len(conn) == 0 {
		panic(fmt.Sprintf("数据链接不可用: %s", conn))
	}

	c, err := gormadapter.NewAdapter(config.Config.DB.Adapter, conn, true) // Your driver and data source.
	if err != nil {
		panic(fmt.Sprintf("NewAdapter 错误: %v", err))
	}

	casbinmodelpath := filepath.Join(config.Root, "config", "rbac_model.conf")
	Enforcer, err = casbin.NewEnforcer(casbinmodelpath, c)
	if err != nil {
		panic(fmt.Sprintf("NewEnforcer 错误: %v", err))
	}

	_ = Enforcer.LoadPolicy()
}
