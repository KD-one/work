package dao

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"test/model"
)

var dB *gorm.DB

func InitDB() {
	db, err := gorm.Open(mysql.Open(viper.GetString("mysql.dsn")))
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}

	var c model.ECUVer
	_ = db.AutoMigrate(&c)

	dB = db
}
