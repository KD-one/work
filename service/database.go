package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"regexp"
	"strings"
	"test/common"
	"test/model"
)

// GetTableData 获取数据库数据,创建新表后(不包括paraauth_开头的表)需要手动在此函数中添加对应数据表的case,并手动定义相应结构体
func GetTableData(c *gin.Context) {
	// 保存数据变量
	var data interface{}
	// 从参数中获取表名
	table := c.Query("table")

	// 使用正则表达式进行paraauth_前缀的验证
	paraauth_prefix := regexp.MustCompile("^paraauth_")

	switch {
	case table == "filemaps":
		var f []model.Filemap
		common.DB.Table(table).Find(&f)
		data = f
	case table == "files":
		var f []model.File
		common.DB.Table(table).Find(&f)
		data = f
	case table == "users":
		var u []model.User
		common.DB.Table(table).Find(&u)
		data = u
	case table == "appauth":
		var a []model.Appauth
		common.DB.Table(table).Find(&a)
		data = a
	case paraauth_prefix.MatchString(table):
		var a []model.Paraauth
		common.DB.Table(table).Find(&a)
		data = a
	}
	if data == nil {
		c.JSON(400, gin.H{
			"msg":  "error",
			"data": "没有找到对应表",
		})
	} else {
		c.JSON(200, gin.H{
			"msg":  "success",
			"data": data,
		})
	}
}

//	CreateTable 传参格式：{
//	   "table": "test",
//	   "fields":[
//	       {
//	           "name": "example_field",
//	           "type": "tinyint",
//	           "is_nullable": false,
//	           "constraints": {
//	               "default": "0"
//	           }
//	       }
//	   ]
//	}
func CreateTable(c *gin.Context) {
	var receivedJson model.DynamicModel
	if err := c.ShouldBindJSON(&receivedJson); err != nil {
		c.JSON(400, gin.H{"error": "错误的json数据！！"})
		return
	}

	// 根据 JSON 数据构造表结构
	var fieldsStr []string
	for _, field := range receivedJson.FieldsInfo {
		// 根据数据库类型和约束条件生成字段定义字符串
		fieldDef := fmt.Sprintf("`%s` %s", field.Name, field.Type)
		if field.Constraints.Length > 0 {
			fieldDef += fmt.Sprintf("(%d)", field.Constraints.Length)
		}
		if !field.IsNullable {
			fieldDef += " NOT NULL"
		}
		if field.Constraints.Default != "" {
			fieldDef += fmt.Sprintf(" DEFAULT '%s'", field.Constraints.Default)
		}
		fieldsStr = append(fieldsStr, fieldDef)
	}

	// 构造创建表的 SQL 语句
	createTableSql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (%s);", receivedJson.TableName, strings.Join(fieldsStr, ", "))

	// 执行 SQL 语句创建表
	common.DB.Exec(createTableSql)

	c.JSON(200, gin.H{"message": fmt.Sprintf("Table `%s` created", receivedJson.TableName)})
}

// InsertOrUpdate 创建或更新数据，创建新表后(仅限不包括paraauth_开头的表)需要手动在此函数中添加对应数据表的case,并手动定义相应结构体
func InsertOrUpdate(c *gin.Context) {

	table := c.Query("table")

	// 使用正则表达式进行paraauth_前缀的验证
	paraauth_prefix := regexp.MustCompile("^paraauth_")

	switch {
	case table == "filemaps":
		type array struct {
			Records []model.Filemap `json:"records"`
		}
		var t array

		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(400, gin.H{"error": "错误的json数据！！"})
			return
		}
		for _, item := range t.Records {
			// 没找到数据就创建
			if res := common.DB.Table(table).Where("id = ?", item.ID).First(&model.Filemap{}); res.RowsAffected == 0 {
				if err := common.DB.Table(table).Create(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("插入记录时出错： %v, %d", err, res.RowsAffected)})
					return
				}
			} else {
				if err := common.DB.Table(table).Updates(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("更新记录时出错： %v", err)})
					return
				}
			}
		}
	case table == "files":
		type array struct {
			Records []model.File `json:"records"`
		}
		var t array

		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(400, gin.H{"error": "错误的json数据！！"})
			return
		}
		for _, item := range t.Records {
			// 没找到数据就创建
			if res := common.DB.Table(table).Where("id = ?", item.Uuid).First(&model.File{}); res.RowsAffected == 0 {
				if err := common.DB.Table(table).Create(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("插入记录时出错： %v, %d", err, res.RowsAffected)})
					return
				}
			} else {
				if err := common.DB.Table(table).Updates(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("更新记录时出错： %v", err)})
					return
				}
			}
		}
	case table == "users":
		type array struct {
			Records []model.User `json:"records"`
		}
		var t array

		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(400, gin.H{"error": "错误的json数据！！"})
			return
		}
		for _, item := range t.Records {
			// 没找到数据就创建
			if res := common.DB.Table(table).Where("id = ?", item.ID).First(&model.User{}); res.RowsAffected == 0 {
				if err := common.DB.Table(table).Create(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("插入记录时出错： %v, %d", err, res.RowsAffected)})
					return
				}
			} else {
				if err := common.DB.Table(table).Updates(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("更新记录时出错： %v", err)})
					return
				}
			}
		}
	case table == "appauth":
		type array struct {
			Records []model.Appauth `json:"records"`
		}
		var t array

		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(400, gin.H{"error": "错误的json数据！！"})
			return
		}
		for _, item := range t.Records {
			// 没找到数据就创建
			if res := common.DB.Table(table).Where("id = ?", item.Id).First(&model.Appauth{}); res.RowsAffected == 0 {
				if err := common.DB.Table(table).Create(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("插入记录时出错： %v, %d", err, res.RowsAffected)})
					return
				}
			} else {
				if err := common.DB.Table(table).Updates(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("更新记录时出错： %v", err)})
					return
				}
			}
		}
	case paraauth_prefix.MatchString(table):
		type array struct {
			Records []model.Paraauth `json:"records"`
		}
		var t array

		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(400, gin.H{"error": "错误的json数据！！"})
			return
		}
		for _, item := range t.Records {
			// 没找到数据就创建
			if res := common.DB.Table(table).Where("id = ?", item.Id).First(&model.Paraauth{}); res.RowsAffected == 0 {
				if err := common.DB.Table(table).Create(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("插入记录时出错： %v, %d", err, res.RowsAffected)})
					return
				}
			} else {
				if err := common.DB.Table(table).Updates(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("更新记录时出错： %v", err)})
					return
				}
			}
		}
	case table == "test":
		type array struct {
			Records []model.Test `json:"records"`
		}
		var t array

		if err := c.ShouldBindJSON(&t); err != nil {
			c.JSON(400, gin.H{"error": "错误的json数据！！"})
			return
		}
		for _, item := range t.Records {
			// 没找到数据就创建
			if res := common.DB.Table(table).Where("id = ?", item.Id).First(&model.Test{}); res.RowsAffected == 0 {
				if err := common.DB.Table(table).Create(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("插入记录时出错： %v, %d", err, res.RowsAffected)})
					return
				}
			} else {
				if err := common.DB.Table(table).Updates(&item).Error; err != nil {
					c.JSON(500, gin.H{"error": fmt.Sprintf("更新记录时出错： %v", err)})
					return
				}
			}
		}
	default:
		c.JSON(400, gin.H{"error": fmt.Sprintf("没找到表: %s", table)})
		return
	}
}

// GetTables 获取数据库所有表的数据
func GetTables(c *gin.Context) {
	// 使用正则解析出数据库名
	databaseName := regexp.MustCompile(`(?<=\/)[^?]+`).FindString(viper.GetString("mysql.dsn"))
	rows, err := common.DB.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = ?;", databaseName).Rows()
	if err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("获取数据库所有表数据时出错： %v", err)})
		return
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("获取数据库所有表数据时出错： %v", err)})
			return
		}
		tableNames = append(tableNames, tableName)
	}

	if err := rows.Err(); err != nil {
		c.JSON(500, gin.H{"error": fmt.Sprintf("获取数据库所有表数据时出错： %v", err)})
		return
	}

	c.JSON(200, gin.H{"data": tableNames})
}
