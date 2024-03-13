package dao

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"regexp"
	"strings"
	"test/common"
	"test/model"
	"time"
)

func CreateOrUpdateFile(fileName, size, t, projectNumber, versionNumber string) error {
	if fileName != "" {
		tmp := strings.Split(fileName, ".")
		filename := tmp[0]
		fileType := tmp[1]
		f3 := common.DB.Where(map[string]interface{}{
			"name": filename,
			"type": fileType,
		}).First(&model.File{})

		//var file model.File
		// 数据库中没有查到此文件纪录，则插入
		if f3.RowsAffected == 0 {
			// 创建文件记录
			common.DB.Model(&model.File{}).Create(map[string]interface{}{
				"name":           filename,
				"type":           fileType,
				"size":           size,
				"created_at":     t,
				"project_number": projectNumber,
				"version_number": versionNumber,
			})
			//log.Printf("上传的文件名：%s", filename)
		} else {
			return errors.New("文件名已存在")
		}
		// 数据库中查到此文件纪录，则更新此文件记录的项目名和版本号
		//file = model.File{
		//	ProjectNumber: projectNumber,
		//	VersionNumber: versionNumber,
		//}
		//common.DB.Model(&model.File{}).Where("name = ? and type = ?", filename, fileType).Updates(&file)
	}
	return nil
}

func GetFileMapByProjectAndVersion(projectName string, versionNumber string) (fileMap *model.Filemap, err error) {
	err = common.DB.Model(&model.Filemap{}).Where("software_version_branch = ? AND software_version_number = ?", projectName, versionNumber).First(&fileMap).Error
	return
}

func UpdateFileMapByProjectAndVersion(fileMap *model.Filemap, projectName string, versionNumber string) error {
	err := common.DB.Model(&model.Filemap{}).Where("software_version_branch = ? AND software_version_number = ?", projectName, versionNumber).Updates(&fileMap).Error
	return err
}

func DBUserLogin(name, password string, userId *uint) error {
	var user model.User
	res := common.DB.Where("name = ?", name).First(&user)
	if res.RowsAffected == 0 {
		return errors.New("DBUserLogin 用户不存在")
	}

	// 如果输入的密码加密完 和数据库中的密码相同，则证明密码输入正确，否则错误
	if mD5([]byte(password)) != user.Password {
		return errors.New("DBUserLogin 密码错误")
	}
	*userId = user.ID
	return nil
}

func DBUserRegister(name, password string, appAuth, paraAuth string, user *model.User) error {
	common.DB.Where("name = ?", name).First(&user)
	if user.ID != 0 {
		return errors.New("DBUserRegister 用户已存在")
	}

	user.Name = name
	user.Password = mD5([]byte(password))
	user.AppAuth = appAuth
	user.ParaAuth = paraAuth
	if common.DB.Create(&user).Error != nil {
		return errors.New("DBUserRegister 创建用户失败")
	}

	return nil
}

func DBUserAddUpdate(adminId uint, user model.User, version int, changelog string) error {
	err := DBTableCheckVer(version)
	if err != nil {
		return err
	}
	var admin model.User
	// u保存待修改用户信息
	var u model.User
	res := common.DB.Model(model.User{}).Where("name = ?", user.Name).First(&u)

	err = common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error
	if err != nil {
		return errors.New("DBUserAddUpdate adminId不存在")
	}
	if admin.UserLevel > u.UserLevel {
		// 如果数据库中不存在此用户，则创建用户
		if res.RowsAffected == 0 {
			user.UserLevel = 1
			if common.DB.Model(model.User{}).Create(&user).Error != nil {
				return errors.New("DBUserAddUpdate 创建用户失败")
			}
		} else { // 如果数据库中存在此用户，则更新用户信息
			if common.DB.Model(model.User{}).Where("id = ?", u.ID).Updates(&user).Error != nil {
				return errors.New("DBUserAddUpdate 更新用户失败")
			}
		}
		return DBTableAddNewVer(admin.Name, changelog)
	} else {
		return errors.New("DBUserAddUpdate 权限不足")
	}
}

func DBUserDelete(adminId uint, userName string, version int, changelog string) error {
	err := DBTableCheckVer(version)
	if err != nil {
		return err
	}
	var admin model.User
	// u保存待删除用户信息
	var u model.User
	res := common.DB.Model(model.User{}).Where("name = ?", userName).First(&u)
	// 判断adminId是否存在,存在则赋值给admin变量
	if common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBUserDelete adminId不存在")
	}

	if res.RowsAffected == 0 {
		return errors.New("DBUserDelete 用户不存在")
	} else {
		if admin.UserLevel > u.UserLevel {
			if common.DB.Model(model.User{}).Delete(&u).Error != nil {
				return errors.New("DBUserDelete 删除失败")
			}
		} else {
			return errors.New("DBUserDelete 权限不足")
		}
	}
	return DBTableAddNewVer(admin.Name, changelog)
}

func DBUserGetTable(adminId uint, users *[]model.User) error {
	var admin model.User
	if common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBUserGetTable adminId不存在")
	}
	if admin.UserLevel < 2 {
		return errors.New("DBUserGetTable 权限不足")
	} else {
		if common.DB.Model(&model.User{}).Find(&users).Error != nil {
			return errors.New("DBUserGetTable 查询失败")
		}
	}
	return nil
}

func DBAppAuthGetUserAuth(adminId uint, appAuth *model.Appauth) error {
	// u保存待查询的用户信息
	var u model.User
	res := common.DB.Model(model.User{}).Where("id = ?", adminId).First(&u)
	if res.RowsAffected == 0 {
		return errors.New("DBAppAuthGetUserAuth 用户不存在")
	} else {
		err := common.DB.Model(&model.Appauth{}).Where("name = ?", u.AppAuth).First(&appAuth).Error
		if err != nil {
			return errors.New("DBAppAuthGetUserAuth appAuth不存在")
		}
	}
	return nil
}

func DBAppAuthGetTable(adminId uint, appAuths *[]model.Appauth) error {
	var admin model.User
	if common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBAppAuthGetTable 用户不存在")
	}
	if admin.UserLevel < 2 {
		return errors.New("DBAppAuthGetTable 权限不足")
	} else {
		if common.DB.Model(&model.Appauth{}).Find(&appAuths).Error != nil {
			return errors.New("DBAppAuthGetTable 查询失败")
		}
	}
	return nil
}

func DBAppAuthAddUpdate(adminId uint, appAuth model.Appauth, version int, changelog string) error {
	err := DBTableCheckVer(version)
	if err != nil {
		return err
	}
	var admin model.User
	res := common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBAppAuthAddUpdate 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBAppAuthAddUpdate 权限不足")
		} else {
			r := common.DB.Model(&model.Appauth{}).Where("name = ?", appAuth.Name).First(&model.Appauth{})
			if r.RowsAffected == 0 {
				if common.DB.Model(&model.Appauth{}).Create(&appAuth).Error != nil {
					return errors.New("DBAppAuthAddUpdate 添加失败")
				}
			} else {
				if common.DB.Model(&model.Appauth{}).Where("name = ?", appAuth.Name).Updates(&appAuth).Error != nil {
					return errors.New("DBAppAuthAddUpdate 更新失败")
				}
			}
		}
	}
	return DBTableAddNewVer(admin.Name, changelog)
}

func DBAppAuthDelete(adminId uint, name string, version int, changelog string) error {
	err := DBTableCheckVer(version)
	if err != nil {
		return err
	}
	var admin model.User
	res := common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBAppAuthDelete 用户不存在")
	} else if admin.UserLevel < 2 {
		return errors.New("DBAppAuthDelete 权限不足")
	} else {
		// 检查是否有用户正在使用这个name
		var count int64
		if common.DB.Model(&model.User{}).
			Where("app_auth = ?", name).
			Count(&count).Error != nil {
			return errors.New("DBAppAuthDelete 查找用户的app_auth字段个数时失败")
		}

		if count > 0 {
			return errors.New("DBAppAuthDelete 无法删除，有用户正在使用此AppAuth名称，请先移除所有关联")
		}

		// 检查并删除AppAuth记录
		var existingAppAuth model.Appauth
		r := common.DB.Model(&model.Appauth{}).Where("name = ?", name).First(&existingAppAuth)
		if r.RowsAffected == 0 {
			return errors.New("DBAppAuthDelete 记录不存在，请检查输入是否正确")
		} else {
			if common.DB.Delete(&existingAppAuth).Error != nil {
				return errors.New("DBAppAuthDelete 删除失败")
			}
		}
	}
	return DBTableAddNewVer(admin.Name, changelog)
}

func DBParaAuthGetUserAuth(adminId uint, paraAuth *model.Paraauth) error {
	// u保存待查询的用户信息
	var u model.User
	res := common.DB.Model(model.User{}).Where("id = ?", adminId).First(&u)
	if res.RowsAffected == 0 {
		return errors.New("DBParaAuthGetUserAuth 用户不存在")
	} else {
		err := common.DB.Model(&model.Paraauth{}).Where("name = ?", u.ParaAuth).First(&paraAuth).Error
		if err != nil {
			return errors.New("DBParaAuthGetUserAuth paraAuth不存在")
		}
	}
	return nil
}

func DBParaAuthGetTable(adminId uint, tableName string, paraAuths *[]model.Paraauth) error {
	var admin model.User
	if common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBParaAuthGetTable adminId不存在")
	}
	if admin.UserLevel < 2 {
		return errors.New("DBParaAuthGetTable 权限不足")
	} else {
		err := common.DB.Table(tableName).Find(&paraAuths).Error
		if err != nil {
			return errors.New("DBParaAuthGetTable paraAuth表不存在")
		}
	}
	return nil
}

func DBParaAuthGetTableList(adminId uint, tableNames *[]string) error {
	var admin model.User
	if common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBParaAuthGetTableList adminId不存在")
	}
	if admin.UserLevel < 2 {
		return errors.New("DBParaAuthGetTableList 权限不足")
	} else {
		// 使用正则解析出数据库名
		databaseName := regexp.MustCompile(`(?<=\/)[^?]+`).FindString(viper.GetString("mysql.dsn"))
		rows, err := common.DB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = ?;", databaseName).Rows()
		if err != nil {
			return errors.New("DBParaAuthGetTableList 获取数据库所有paraauth_表数据时出错")
		}
		defer rows.Close()

		for rows.Next() {
			var tableName string
			if err := rows.Scan(&tableName); err != nil {
				return errors.New("DBParaAuthGetTableList 获取数据库所有paraauth_表数据时出错")
			}
			if strings.HasPrefix(tableName, "paraauth_") {
				*tableNames = append(*tableNames, tableName)
			}
		}

		if err = rows.Err(); err != nil {
			return errors.New("DBParaAuthGetTableList 获取数据库所有paraauth_表数据时出错")
		}
	}
	return nil
}

func DBParaAuthAddUpdateTable(adminId uint, paraTableName string, paraAuth []model.Paraauth, version int, changelog string) error {
	err := DBTableCheckVer(version)
	if err != nil {
		return err
	}
	var admin model.User
	res := common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBParaAuthAddUpdateTable 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBParaAuthAddUpdateTable 权限不足")
		} else {
			// 检查并处理表
			if checkTableExists(paraTableName) {
				// 表存在，先删除再重建
				if err := common.DB.Migrator().DropTable(paraTableName); err != nil {
					return fmt.Errorf("DBParaAuthAddUpdateTable 删除表 %s 时出错: %s", paraTableName, err.Error())
				}
			}

			// 创建新表
			createTableSQL := fmt.Sprintf("CREATE TABLE %s (id int unsigned NOT NULL AUTO_INCREMENT,para_name varchar(255) UNIQUE NOT NULL,min_value float NOT NULL,max_value float NOT NULL,change_enable tinyint(1) NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB;", paraTableName)

			if err := common.DB.Exec(createTableSQL).Error; err != nil {
				return fmt.Errorf("DBParaAuthAddUpdateTable 创建表 %s 时出错: %s", paraTableName, err.Error())
			}
			// 将 Paraauth 结构体数组数据插入新表
			for _, auth := range paraAuth {
				if err := common.DB.Table(paraTableName).Create(&auth).Error; err != nil {
					return fmt.Errorf("DBParaAuthAddUpdateTable 插入数据到 %s 表时出错: %s", paraTableName, err.Error())
				}
			}
		}
	}
	return DBTableAddNewVer(admin.Name, changelog)
}

func DBParaAuthDeleteTable(adminId uint, paraTableName string, version int, changelog string) error {
	err := DBTableCheckVer(version)
	if err != nil {
		return err
	}
	var admin model.User
	res := common.DB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBParaAuthDeleteTable 用户不存在")
	} else if admin.UserLevel < 2 {
		return errors.New("DBParaAuthDeleteTable 权限不足")
	} else {
		if checkTableExists(paraTableName) {
			// 检查是否有用户正在使用这个para_auth表名
			var count int64
			if common.DB.Model(&model.User{}).
				Where("para_auth = ?", paraTableName).
				Count(&count).Error != nil {
				return errors.New("DBParaAuthDeleteTable 查找用户的para_auth字段个数时失败")
			}

			if count > 0 {
				return errors.New("DBParaAuthDeleteTable 无法删除，有用户正在使用此ParaAuth名称，请先移除所有关联")
			}
			if err := common.DB.Migrator().DropTable(paraTableName); err != nil {
				return fmt.Errorf("DBParaAuthDeleteTable 删除表 %s 时出错: %s", paraTableName, err.Error())
			}
		} else {
			return errors.New("DBParaAuthDeleteTable 表不存在，请检查名称")
		}
	}
	return DBTableAddNewVer(admin.Name, changelog)
}

func DBTableGetLatestVer(version *int) error {
	var maxRecord model.Tablever
	if common.DB.Model(&model.Tablever{}).Order("ver desc").First(&maxRecord).Error != nil {
		return errors.New("获取最新版本记录时出错")
	}
	*version = maxRecord.Ver
	return nil
}

func DBTableAddNewVer(name, changelog string) error {
	record := model.Tablever{
		User:      name,
		ChangeLog: changelog,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	err := common.DB.Create(&record).Error
	if err != nil {
		return errors.New("DBTableAddNewVer 创建新的版本记录时出错")
	}
	return nil
}

func checkTableExists(tableName string) bool {
	// 获取 migrator 接口
	migrator := common.DB.Migrator()

	// 检查表是否存在
	exists := migrator.HasTable(tableName)

	return exists
}

func DBTableCheckVer(version int) error {
	var latesversion int
	if err := DBTableGetLatestVer(&latesversion); err != nil {
		return err
	}
	if version != latesversion {
		fmt.Println(latesversion)
		fmt.Println(version)
		var versionRecord model.Tablever
		if common.DB.Model(model.Tablever{}).Where("ver = ?", latesversion).First(&versionRecord).Error != nil {
			return errors.New("DBTableCheckVer 获取最新版本记录时出错")
		}
		return errors.New(fmt.Sprintf("DBTableCheckVer 数据已被%s在%s时间修改，修改内容为%s", versionRecord.User, versionRecord.CreatedAt, versionRecord.ChangeLog))
	}
	return nil
}

//  系统使用函数，慎用！！！！！！！！！！！！！！！！！！！！！！！！！

func FindUserTable(users *[]model.User) error {
	if common.DB.Model(&model.User{}).Find(&users).Error != nil {
		return errors.New("DBUserGetTable 查询失败")
	}
	return nil
}
