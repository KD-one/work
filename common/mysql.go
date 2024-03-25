package common

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"test/model"
)

var DB *gorm.DB

func InitDB() {
	//host := "127.0.0.1"
	//port := "3306"
	//database := "user"
	//username := "root"
	//password := "root"
	//charset := "utf8"
	//dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
	//	username,
	//	password,
	//	host,
	//	port,
	//	database,
	//	charset)
	db, err := gorm.Open(mysql.Open(viper.GetString("mysql.dsn")))
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}

	//var user model.User
	//db.AutoMigrate(&user)

	//var f model.File
	//_ = db.AutoMigrate(&f)

	//var f model.Filemap
	//_ = db.AutoMigrate(&f)

	//var t model.Tablever
	//_ = db.AutoMigrate(&t)

	var c model.Appauth
	_ = db.AutoMigrate(&c)

	DB = db
}
