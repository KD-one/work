package dao

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"test/model"
	"time"
)

// -------------------------------------------- // --------------------------------------------------- // ------------------------------------------------ //

func DBUserLogin(name, password string, userId *uint) error {
	var user model.User
	res := dB.Where("name = ?", name).First(&user)
	if res.RowsAffected == 0 {
		return errors.New("DBUserLogin 用户不存在")
	}

	// 如果输入的密码加密完 和数据库中的密码相同，则证明密码输入正确，否则错误
	if password != user.Password {
		return errors.New("DBUserLogin 密码错误")
	}
	*userId = user.ID
	return nil
}

func DBUserGetTable(adminId uint, users *[]model.User) error {
	var admin model.User
	if dB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBUserGetTable adminId不存在")
	}
	if admin.UserLevel < 2 {
		return errors.New("DBUserGetTable 权限不足")
	} else {
		if dB.Model(&model.User{}).Find(&users).Error != nil {
			return errors.New("DBUserGetTable 查询失败")
		}
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
	res := dB.Model(model.User{}).Where("name = ?", user.Name).First(&u)

	err = dB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error
	if err != nil {
		return errors.New("DBUserAddUpdate adminId不存在")
	}
	if admin.UserLevel > u.UserLevel {
		// 如果数据库中不存在此用户，则创建用户
		if res.RowsAffected == 0 {
			user.UserLevel = 1
			if dB.Model(model.User{}).Create(&user).Error != nil {
				return errors.New("DBUserAddUpdate 创建用户失败")
			}
		} else { // 如果数据库中存在此用户，则更新用户信息
			if dB.Model(model.User{}).Where("id = ?", u.ID).Updates(&user).Error != nil {
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
	res := dB.Model(model.User{}).Where("name = ?", userName).First(&u)
	// 判断adminId是否存在,存在则赋值给admin变量
	if dB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBUserDelete adminId不存在")
	}

	if res.RowsAffected == 0 {
		return errors.New("DBUserDelete 用户不存在")
	} else {
		if admin.UserLevel > u.UserLevel {
			if dB.Model(model.User{}).Delete(&u).Error != nil {
				return errors.New("DBUserDelete 删除失败")
			}
		} else {
			return errors.New("DBUserDelete 权限不足")
		}
	}
	return DBTableAddNewVer(admin.Name, changelog)
}

func DBAppAuthGetUserAuth(adminId uint, appAuth *model.Appauth) error {
	// u保存待查询的用户信息
	var u model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&u)
	if res.RowsAffected == 0 {
		return errors.New("DBAppAuthGetUserAuth 用户不存在")
	} else {
		if u.UserLevel < 2 {
			return errors.New("DBAppAuthGetUserAuth 权限不足")
		} else {
			err := dB.Model(&model.Appauth{}).Where("auth_name = ?", u.AppAuth).First(&appAuth).Error
			if err != nil {
				return errors.New("DBAppAuthGetUserAuth appAuth不存在")
			}
		}
	}
	return nil
}

func DBFindAppAuthByName(adminId uint, authName string, appAuth *model.Appauth) error {
	var admin model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBFindAppAuthByName 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBFindAppAuthByName 权限不足")
		} else {
			err := dB.Model(&model.Appauth{}).Where("auth_name = ?", authName).First(&appAuth).Error
			if err != nil {
				return errors.New("DBFindAppAuthByName appAuth不存在")
			}
		}
	}
	return nil
}

func DBAppAuthGetTable(adminId uint, appAuths *[]model.Appauth) error {
	var admin model.User
	if dB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBAppAuthGetTable 用户不存在")
	}
	if admin.UserLevel < 2 {
		return errors.New("DBAppAuthGetTable 权限不足")
	} else {
		if dB.Model(&model.Appauth{}).Find(&appAuths).Error != nil {
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
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBAppAuthAddUpdate 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBAppAuthAddUpdate 权限不足")
		} else {
			r := dB.Model(&model.Appauth{}).Where("auth_name = ?", appAuth.AuthName).First(&model.Appauth{})
			if r.RowsAffected == 0 {
				if dB.Model(&model.Appauth{}).Create(&appAuth).Error != nil {
					return errors.New("DBAppAuthAddUpdate 添加失败")
				}
			} else {
				if dB.Model(&model.Appauth{}).Where("auth_name = ?", appAuth.AuthName).Updates(&appAuth).Error != nil {
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
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBAppAuthDelete 用户不存在")
	} else if admin.UserLevel < 2 {
		return errors.New("DBAppAuthDelete 权限不足")
	} else {
		// 检查是否有用户正在使用这个name
		var count int64
		if dB.Model(&model.User{}).
			Where("app_auth = ?", name).
			Count(&count).Error != nil {
			return errors.New("DBAppAuthDelete 查找用户的app_auth字段个数时失败")
		}

		if count > 0 {
			return errors.New("DBAppAuthDelete 无法删除，有用户正在使用此AppAuth名称，请先移除所有关联")
		}

		// 检查并删除AppAuth记录
		var existingAppAuth model.Appauth
		r := dB.Model(&model.Appauth{}).Where("auth_name = ?", name).First(&existingAppAuth)
		if r.RowsAffected == 0 {
			return errors.New("DBAppAuthDelete 记录不存在，请检查输入是否正确")
		} else {
			if dB.Delete(&existingAppAuth).Error != nil {
				return errors.New("DBAppAuthDelete 删除失败")
			}
		}
	}
	return DBTableAddNewVer(admin.Name, changelog)
}

func DBParaAuthGetUserAuth(adminId uint, paraAuth *model.Paraauth) error {
	// u保存待查询的用户信息
	var u model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&u)
	if res.RowsAffected == 0 {
		return errors.New("DBParaAuthGetUserAuth 用户不存在")
	} else {
		err := dB.Model(&model.Paraauth{}).Where("para_name = ?", u.ParaAuth).First(&paraAuth).Error
		if err != nil {
			return errors.New("DBParaAuthGetUserAuth paraAuth不存在")
		}
	}
	return nil
}

func DBParaAuthGetTable(adminId uint, tableName string, paraAuths *[]model.Paraauth) error {
	var admin model.User
	if dB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBParaAuthGetTable adminId不存在")
	}
	if admin.UserLevel < 2 {
		return errors.New("DBParaAuthGetTable 权限不足")
	} else {
		err := dB.Table(tableName).Find(&paraAuths).Error
		if err != nil {
			return errors.New("DBParaAuthGetTable paraAuth表不存在")
		}
	}
	return nil
}

func DBParaAuthGetTableList(adminId uint, tableNames *[]string) error {
	var admin model.User
	if dB.Model(model.User{}).Where("id = ?", adminId).First(&admin).Error != nil {
		return errors.New("DBParaAuthGetTableList adminId不存在")
	}
	if admin.UserLevel < 2 {
		return errors.New("DBParaAuthGetTableList 权限不足")
	} else {
		// 使用正则解析出数据库名
		databaseName := viper.GetString("mysql.database")
		rows, err := dB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = ?;", databaseName).Rows()
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
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBParaAuthAddUpdateTable 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBParaAuthAddUpdateTable 权限不足")
		} else {
			// 检查并处理表
			if checkTableExists(paraTableName) {
				// 表存在，先删除再重建
				if err := dB.Migrator().DropTable(paraTableName); err != nil {
					return fmt.Errorf("DBParaAuthAddUpdateTable 删除表 %s 时出错: %s", paraTableName, err.Error())
				}
			}

			// 创建新表
			createTableSQL := fmt.Sprintf("CREATE TABLE %s (id int unsigned NOT NULL AUTO_INCREMENT,para_name varchar(255) UNIQUE NOT NULL,min_value float NOT NULL,max_value float NOT NULL,change_enable tinyint(1) NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB;", paraTableName)

			if err := dB.Exec(createTableSQL).Error; err != nil {
				return fmt.Errorf("DBParaAuthAddUpdateTable 创建表 %s 时出错: %s", paraTableName, err.Error())
			}
			// 将 Paraauth 结构体数组数据插入新表
			for _, auth := range paraAuth {
				if err := dB.Table(paraTableName).Create(&auth).Error; err != nil {
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
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBParaAuthDeleteTable 用户不存在")
	} else if admin.UserLevel < 2 {
		return errors.New("DBParaAuthDeleteTable 权限不足")
	} else {
		if checkTableExists(paraTableName) {
			// 检查是否有用户正在使用这个para_auth表名
			var count int64
			if dB.Model(&model.User{}).
				Where("para_auth = ?", paraTableName).
				Count(&count).Error != nil {
				return errors.New("DBParaAuthDeleteTable 查找用户的para_auth字段个数时失败")
			}

			if count > 0 {
				return errors.New("DBParaAuthDeleteTable 无法删除，有用户正在使用此ParaAuth名称，请先移除所有关联")
			}
			if err := dB.Migrator().DropTable(paraTableName); err != nil {
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
	if dB.Model(&model.Tablever{}).Order("ver desc").First(&maxRecord).Error != nil {
		return errors.New("DBTableGetLatestVer 获取最新版本记录时出错")
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
	err := dB.Create(&record).Error
	if err != nil {
		return errors.New("DBTableAddNewVer 创建新的版本记录时出错")
	}
	return nil
}

func checkTableExists(tableName string) bool {
	// 获取 migrator 接口
	migrator := dB.Migrator()

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
		fmt.Println("latesversion: ", latesversion)
		fmt.Println("version: ", version)
		var versionRecord model.Tablever
		if dB.Model(model.Tablever{}).Where("ver = ?", latesversion).First(&versionRecord).Error != nil {
			return errors.New("DBTableCheckVer 获取最新版本记录时出错")
		}
		return errors.New(fmt.Sprintf("DBTableCheckVer 数据已被%s在%s时间修改，修改内容为%s", versionRecord.User, versionRecord.CreatedAt, versionRecord.ChangeLog))
	}
	return nil
}

func DBGetLimitVersionTable(adminId uint, tableVersion *[]model.Tablever) error {
	var admin model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBGetLimitVersionTable 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBGetLimitVersionTable 权限不足")
		} else {
			if dB.Model(model.Tablever{}).Order("ver desc").Limit(viper.GetInt("versionRecord.limit")).Find(&tableVersion).Error != nil {
				return errors.New("DBGetLimitVersionTable 获取版本记录时出错")
			}
		}
	}
	return nil
}

func DBGetLimitClientLogTable(adminId uint, clientLog *[]model.ClientLog) error {
	var admin model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBGetLimitClientLogTable 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBGetLimitClientLogTable 权限不足")
		} else {
			if dB.Model(model.ClientLog{}).Order("id desc").Limit(viper.GetInt("versionRecord.limit")).Find(&clientLog).Error != nil {
				return errors.New("DBGetLimitClientLogTable 获取版本记录时出错")
			}
		}
	}
	return nil
}

func DBCheckECUFileMapRecordExists(branch, version uint) int64 {
	m := dB.Where(map[string]interface{}{
		"branch":  branch,
		"version": version,
	}).First(&model.EcuFileMap{})
	return m.RowsAffected
}

func DBCreateECUFileMapRecord(adminId uint, ecuFileMap model.EcuFileMap) error {
	var admin model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBCreateECUFileMapRecord 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBCreateECUFileMapRecord 权限不足")
		} else {
			err := dB.Create(&ecuFileMap).Error
			if err != nil {
				return errors.New("DBCreateECUFileMapRecord 插入记录时出错")
			}
		}
	}
	return nil
}

func DBWhereBuildFileFindRecord(buildFile string, ecuFileMap *model.EcuFileMap) error {
	if dB.Model(&model.EcuFileMap{}).Where("build_file = ?", buildFile).First(&ecuFileMap).Error != nil {
		return errors.New("DBWhereBuildFileFindRecord 查找记录时出错")
	}
	return nil
}

func DBFindECUFileMapRecord(adminId, branch, version uint, ecuFileMap *model.EcuFileMap) error {
	var admin model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("DBFindECUFileMapRecord 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("DBFindECUFileMapRecord 权限不足")
		} else {
			err := dB.Model(&model.EcuFileMap{}).Where("branch = ? AND version = ?", branch, version).First(&ecuFileMap).Error
			if err != nil {
				return errors.New(fmt.Sprintf("DBFindECUFileMapRecord 查找记录时出错 err:%s", err.Error()))
			}
		}
	}
	return nil
}

func DBFindECUFileMapRecordNoAuth(branch, version uint, ecuFileMap *model.EcuFileMap) error {
	err := dB.Model(&model.EcuFileMap{}).Where("branch = ? AND version = ?", branch, version).First(&ecuFileMap).Error
	if err != nil {
		return errors.New(fmt.Sprintf("DBFindECUFileMapRecord 查找记录时出错 err:%s", err.Error()))
	}
	return nil
}

func DBCheckGtCurrentVersion(branch, version uint, ecuFileMap *[]model.EcuFileMap) error {
	err := dB.Model(&model.EcuFileMap{}).Where("branch = ? AND version > ?", branch, version).Find(&ecuFileMap).Error
	if err != nil {
		return errors.New("DBCheckGtCurrentVersion 检查版本时出错")
	}
	return nil
}

func DBClientLogAdd(cLog model.ClientLog) error {
	dB.Begin()
	err := dB.Create(&cLog).Error
	if err != nil {
		return errors.New("DBClientLogAdd 插入数据时出错")
		dB.Rollback()
	}
	dB.Commit()
	return nil
}

func ECUProjectAddChange(adminId uint, ecuProject model.ECUProjectList) error {
	var admin model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("ECUProjectAddChange 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("ECUProjectAddChange 权限不足")
		} else {
			r := dB.Model(&model.ECUProjectList{}).Where("project_code = ?", ecuProject.ProjectCode).First(&model.ECUProjectList{})
			if r.RowsAffected == 0 {
				err := dB.Create(&ecuProject).Error
				if err != nil {
					return errors.New("ECUProjectAddChange 插入数据时出错")
				}
			} else {
				err := dB.Model(&model.ECUProjectList{}).Where("project_code = ?", ecuProject.ProjectCode).Updates(&ecuProject).Error
				if err != nil {
					return errors.New("ECUProjectAddChange 更新数据时出错")
				}
			}
		}
	}
	return nil
}

func DBCheckBranchAndVersionRepeat(branch, version uint) bool {
	row := dB.Model(&model.ECUVer{}).Where("sw_branch = ? AND sw_version = ?", branch, version).First(&model.ECUVer{})
	if row.RowsAffected != 0 {
		return true
	}
	return false
}
func DBFindECUSoftwareRecord(branch, version uint, ecuVer *model.ECUVer) error {
	if dB.Model(&model.ECUVer{}).Where("sw_branch = ? AND sw_version = ?", branch, version).First(&ecuVer).Error != nil {
		return errors.New("DBFindECUSoftwareRecord 查找记录时出错")
	}
	return nil
}

func DBWhereSWBranchFindECUProjectListRecord(branch uint, ecuProjectList *model.ECUProjectList) error {
	if dB.Model(&model.ECUProjectList{}).Where("software_branch = ?", branch).First(&ecuProjectList).Error != nil {
		return errors.New("DBWhereSWBranchFindECUProjectListRecord 查找记录时出错")
	}
	return nil
}

func ECUSoftwareVersionAdd(adminId uint, ecuVer model.ECUVer) error {
	var admin model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("ECUSoftwareVersionAdd 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("ECUSoftwareVersionAdd 权限不足")
		} else {
			if dB.Create(&ecuVer).Error != nil {
				return errors.New("ECUSoftwareVersionAdd 插入数据时出错")
			}
		}
	}
	return nil
}

func ECUSoftwareVersionChange(adminId uint, ecuVer model.ECUVer) error {
	var admin model.User
	res := dB.Model(model.User{}).Where("id = ?", adminId).First(&admin)
	if res.RowsAffected == 0 {
		return errors.New("ECUSoftwareVersionChange 用户不存在")
	} else {
		if admin.UserLevel < 2 {
			return errors.New("ECUSoftwareVersionChange 权限不足")
		} else {
			if dB.Model(&model.ECUVer{}).Where("sw_branch = ? AND sw_version = ?", ecuVer.SWBranch, ecuVer.SWVersion).Updates(&ecuVer).Error != nil {
				return errors.New("ECUSoftwareVersionChange 更新数据时出错")
			}
		}
	}
	return nil
}

func DBValidBranch(branch uint) int64 {
	res := dB.Model(&model.ECUVer{}).Where("sw_branch = ?", branch).First(&model.ECUVer{})
	return res.RowsAffected
}

func DBFindLatestBranchRecord(branch uint, ecuVer *model.ECUVer) error {
	if dB.Model(&model.ECUVer{}).Where("sw_branch = ?", branch).Order("id desc").First(&ecuVer).Error != nil {
		return errors.New("DBCheckSemiFinishedProducts 查询失败")
	}
	return nil
}

//  系统使用函数，慎用！！！！！！！！！！！！！！！！！！！！！！！！！

func FindUserTable(users *[]model.User) error {
	if dB.Model(&model.User{}).Find(&users).Error != nil {
		return errors.New("FindUserTable 查询失败")
	}
	return nil
}

func FindInstruction(instructionId uint, instruction *model.Instruction) error {
	if dB.Model(&model.Instruction{}).Where("id = ?", instructionId).First(&instruction).RowsAffected == 0 {
		return errors.New("FindInstructionTable 没找到记录")
	}
	return nil
}

func UpdateInstructionRecord(instructionId uint, instruction model.Instruction) error {
	if dB.Model(&model.Instruction{}).Where("id = ?", instructionId).Updates(&instruction).Error != nil {
		return errors.New("UpdateInstructionRecord 更新数据时出错")
	}
	return nil
}

func InsertInstructionReturnId(adminId uint, instruction model.Instruction) (uint, error) {
	var admin model.User
	var ins model.Instruction
	res := dB.Model(&model.User{}).Where("id = ?", adminId).First(&admin)

	if admin.UserLevel < 2 {
		return 0, errors.New("InsertInstructionReturnId 权限不足")
	} else {
		if res.RowsAffected == 0 {
			return 0, errors.New("InsertInstructionReturnId 用户不存在")
		} else {
			err := dB.Model(&model.Instruction{}).Create(&instruction).Error
			if err != nil {
				return 0, errors.New("InsertInstructionReturnId 插入数据时出错")
			}
			err = dB.Model(&model.Instruction{}).Order("id desc").First(&ins).Error
			if err != nil {
				return 0, errors.New("InsertInstructionReturnId 获取id时出错")
			}
			return ins.Id, nil
		}
	}
}
func FindUserName(adminId uint) string {
	var admin model.User
	dB.Model(&model.User{}).Where("id = ?", adminId).First(&admin)
	return admin.Name
}
