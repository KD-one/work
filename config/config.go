package config

import (
	"github.com/spf13/viper"
	"test/dao"
	"test/model"
)

func viperInit() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("conf")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

func Init() {
	if err := viperInit(); err != nil {
		panic(err)
	}
}

func InitUserList(userList *[]model.UserList) error {
	var users []model.User
	err := dao.FindUserTable(&users)
	if err != nil {
		return err
	}
	for _, user := range users {
		*userList = append(*userList, model.UserList{
			Name:     user.Name,
			AppAuth:  user.AppAuth,
			ParaAuth: user.ParaAuth,
		})
	}
	return nil
}
