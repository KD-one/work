package common

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	//_ "github.com/go-sql-driver/mysql"
)

var CasBin *casbin.Enforcer

func InitCasbinDB() *casbin.Enforcer {
	// 从数据库加载策略
	host := "127.0.0.1"
	port := "3306"
	database := "test"
	username := "root"
	password := "root"
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s",
		username,
		password,
		host,
		port,
		database)
	adapter, _ := gormadapter.NewAdapter("mysql", dsn)
	CasBin, _ = casbin.NewEnforcer("config/model.conf", adapter)
	//CasBin.AddFunction("ParamsMatch", ParamsMatchFunc)
	err := CasBin.LoadPolicy()
	if err != nil {
		fmt.Printf("LoadPolicy Error! : %s", err)
		return nil
	}
	return CasBin
}
