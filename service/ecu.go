package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"strconv"
	"test/common"
	"test/dao"
	"test/model"
	"test/serializer"
)

type ECUProjectRequestModel struct {
	ChangeInitiator    string                `gorm:"not null" json:"ChangeInitiator" form:"ChangeInitiator"`
	ChangeInitiateTime string                `gorm:"not null" json:"ChangeInitiateTime" form:"ChangeInitiateTime"`
	ChangeCause        string                `gorm:"not null" json:"ChangeCause" form:"ChangeCause"`
	ChangeReq          string                `gorm:"not null" json:"ChangeReq" form:"ChangeReq"`
	ChangeAttached     *multipart.FileHeader `gorm:"not null" json:"ChangeAttached" form:"ChangeAttached"`
	ChangeApplyRange   string                `gorm:"not null" json:"ChangeApplyRange" form:"ChangeApplyRange"`
	SWModifier         string                `gorm:"not null" json:"SWModifier" form:"SWModifier"`
	SWFinishTime       string                `gorm:"not null" json:"SWFinishTime" form:"SWFinishTime"`
	SWLog              string                `gorm:"not null" json:"SWLog" form:"SWLog"`
	SWBuildFile        *multipart.FileHeader `gorm:"not null" json:"SWBuildFile" form:"SWBuildFile"`
	SWA2LFile          *multipart.FileHeader `gorm:"not null" json:"SWA2LFile" form:"SWA2LFile"`
	SWDLLFile          *multipart.FileHeader `gorm:"not null" json:"SWDLLFile" form:"SWDLLFile"`
	SWBranch           string                `gorm:"not null" json:"SWBranch" form:"SWBranch"`
	SWVersion          string                `gorm:"not null" json:"SWVersion" form:"SWVersion"`
	SWCalMain          string                `gorm:"not null" json:"SWCalMain" form:"SWCalMain"`
	SWCalSub           string                `gorm:"not null" json:"SWCalSub" form:"SWCalSub"`
	SWLevel            string                `gorm:"not null" json:"SWLevel" form:"SWLevel"`
	HILTester          string                `gorm:"not null" json:"HILTester" form:"HILTester"`
	HILFinishTime      string                `gorm:"not null" json:"HILFinishTime" form:"HILFinishTime"`
	HILResult          string                `gorm:"not null" json:"HILResult" form:"HILResult"`
	SysVerifier        string                `gorm:"not null" json:"SysVerifier" form:"SysVerifier"`
	SysVerifyTime      string                `gorm:"not null" json:"SysVerifyTime" form:"SysVerifyTime"`
	SysVerifyResult    string                `gorm:"not null" json:"SysVerifyResult" form:"SysVerifyResult"`
	SysVerifyAttached  *multipart.FileHeader `gorm:"not null" json:"SysVerifyAttached" form:"SysVerifyAttached"`
	//SWLevel            uint   `gorm:"not null" json:"SWLevel" form:"SWLevel"`
	SWLogClient   string                `gorm:"not null" json:"SWLogClient" form:"SWLogClient"`
	SWLogClientEN string                `gorm:"not null" json:"SWLogClientEN" form:"SWLogClientEN"`
	CALFile       *multipart.FileHeader `gorm:"not null" json:"CALFile" form:"CALFile"`
	Comment       string                `gorm:"not null" json:"Comment" form:"Comment"`
}

func ECUProjectAddChange(c *gin.Context) {
	var data model.ECUProjectList
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	common.WriteLog(adminId, fmt.Sprintf("ECU项目增加或更新   项目名称：%s   项目编号：%s   软件编号：%d   软件名称：%s", data.ProjectName, data.ProjectCode, data.SoftwareBranch, data.SoftwareName))

	if data.ProjectName == "" || data.ProjectCode == "" || data.SoftwareBranch == 0 || data.SoftwareName == "" {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数不能为空",
		})
		return
	}

	err := dao.ECUProjectAddChange(adminId, data)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}

	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
	})

}

func VerRecordAdd(c *gin.Context) {
	var data ECUProjectRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	common.WriteLog(adminId, fmt.Sprintf("增加版本记录   变更发起人：%s   变更发起时间：%s   变更需求：%s   软件分支号：%s", data.ChangeInitiator, data.ChangeInitiateTime, data.ChangeReq, data.SWBranch))
	if data.ChangeInitiator == "" || data.ChangeInitiateTime == "" || data.ChangeReq == "" || data.SWBranch == "" {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "必选参数不能为空",
		})
		return
	}
	branch, _ := strconv.ParseUint(data.SWBranch, 10, 64)
	version, err := strconv.ParseUint(data.SWVersion, 10, 64)
	if err != nil {
		if data.SWVersion == "" {
			version = 0
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "版本号格式错误",
			})
			return
		}
	}
	calmain, err := strconv.ParseUint(data.SWCalMain, 10, 64)
	if err != nil {
		if data.SWCalMain == "" {
			calmain = 0
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "版本号格式错误",
			})
			return
		}
	}
	calsub, err := strconv.ParseUint(data.SWCalSub, 10, 64)
	if err != nil {
		if data.SWCalSub == "" {
			calsub = 0
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "版本号格式错误",
			})
			return
		}
	}

	// 半成品检查
	var ecuVer model.ECUVer
	if dao.DBValidBranch(uint(branch)) != 0 {
		err = dao.DBFindLatestBranchRecord(uint(branch), &ecuVer)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
		if ecuVer.SWVersion == 0 {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "当前分支存在半成品，无法继续添加",
			})
			return
		}
	}

	if version != 0 {
		if dao.DBCheckBranchAndVersionRepeat(uint(branch), uint(version)) {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "记录已存在，只能修改",
			})
			return
		}
		if uint(version) <= ecuVer.SWVersion {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "当前版本号小于最新的版本号，无法继续添加",
			})
			return
		}
	}
	var output string
	if data.SWBuildFile != nil {
		if data.SWA2LFile == nil || data.SWVersion == "" || data.SWModifier == "" || data.SWFinishTime == "" || data.SWLog == "" {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "SWA2LFile  SWVersion  SWModifier  SWFinishTime  SWLog 这些参数不能为空",
			})
			return
		}
		var epl model.ECUProjectList
		err := dao.DBWhereSWBranchFindECUProjectListRecord(uint(branch), &epl)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
		if data.SWBuildFile.Filename != epl.SoftwareName {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "软件名称不一致",
			})
			return
		}
		buildFilePath := "upload/" + data.SWBuildFile.Filename
		a2lFilePath := "upload/" + data.SWA2LFile.Filename

		_, err = os.Stat(buildFilePath)
		if err == nil {
			err = os.Remove(buildFilePath)
			if err != nil {
				c.JSON(400, serializer.Response{
					Code: 400,
					Msg:  "服务器删除文件时失败",
				})
				return
			}
		}

		err = c.SaveUploadedFile(data.SWBuildFile, buildFilePath)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "服务器写入文件时失败",
			})
			return
		} // 保存文件

		_, err = os.Stat(a2lFilePath)
		if err == nil {
			err = os.Remove(a2lFilePath)
			if err != nil {
				c.JSON(400, serializer.Response{
					Code: 400,
					Msg:  "服务器删除文件时失败",
				})
				return
			}
		}
		err = c.SaveUploadedFile(data.SWA2LFile, a2lFilePath)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "服务器写入文件时失败",
			})
			return
		}
		//}

		files := []string{buildFilePath, a2lFilePath}
		if data.SWDLLFile != nil {
			files = append(files, "upload/"+data.SWDLLFile.Filename)
		}
		output = "ECUSoftware/" + data.SWBuildFile.Filename[:len(data.SWBuildFile.Filename)-4] + ".zip"

		// 判断upload文件夹中是否存在指定的zip文件
		_, err = os.Stat(output)
		if err != nil {
			// 文件不存在，则将文件压缩打包
			if err = common.ZipFiles(output, files); err != nil {
				c.JSON(400, serializer.Response{
					Code: 400,
					Msg:  "压缩文件时失败",
				})
				return
			}
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  fmt.Sprintf("%s 文件已存在", output),
			})
			return
		}
	}

	changeAttached := ""
	sysVerifyAttached := ""
	cALFile := ""
	if data.ChangeAttached == nil {
		changeAttached = ""
	} else {
		changeAttached = data.ChangeAttached.Filename
	}
	if data.SysVerifyAttached == nil {
		sysVerifyAttached = ""
	} else {
		sysVerifyAttached = data.SysVerifyAttached.Filename
	}
	if data.CALFile == nil {
		cALFile = ""
	} else {
		cALFile = data.CALFile.Filename
	}
	if output != "" {
		output = output[12:]
	}
	ecusoftwareversion := model.ECUVer{
		ChangeInitiator:    data.ChangeInitiator,
		ChangeInitiateTime: data.ChangeInitiateTime,
		ChangeCause:        data.ChangeCause,
		ChangeReq:          data.ChangeReq,
		ChangeAttached:     changeAttached,
		ChangeApplyRange:   data.ChangeApplyRange,
		SWModifier:         data.SWModifier,
		SWFinishTime:       data.SWFinishTime,
		SWLog:              data.SWLog,
		SWBuildFile:        output,
		SWBranch:           uint(branch),
		SWVersion:          uint(version),
		SWCalMain:          uint(calmain),
		SWCalSub:           uint(calsub),
		HILTester:          data.HILTester,
		HILFinishTime:      data.HILFinishTime,
		HILResult:          data.HILResult,
		SysVerifier:        data.SysVerifier,
		SysVerifyTime:      data.SysVerifyTime,
		SysVerifyResult:    data.SysVerifyResult,
		SysVerifyAttached:  sysVerifyAttached,
		SWLevel:            4,
		SWLogClient:        data.SWLogClient,
		SWLogClientEN:      data.SWLogClientEN,
		CALFile:            cALFile,
		Comment:            data.Comment,
	}
	err = dao.ECUSoftwareVersionAdd(adminId, ecusoftwareversion)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
	})
}

func VerRecordChange(c *gin.Context) {
	var data ECUProjectRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	common.WriteLog(adminId, fmt.Sprintf("改变版本记录   软件分支号：%s  软件版本号：%s", data.SWBranch, data.SWVersion))

	if data.SWBranch == "" || data.SWVersion == "" {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "必须指定项目号和分支号",
		})
		return
	}

	branch, _ := strconv.ParseUint(data.SWBranch, 10, 64)
	version, err := strconv.ParseUint(data.SWVersion, 10, 64)
	if err != nil {
		if data.SWVersion == "" {
			version = 0
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "版本号格式错误",
			})
			return
		}
	}
	calmain, err := strconv.ParseUint(data.SWCalMain, 10, 64)
	if err != nil {
		if data.SWCalMain == "" {
			calmain = 0
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "版本号格式错误",
			})
			return
		}
	}
	calsub, err := strconv.ParseUint(data.SWCalSub, 10, 64)
	if err != nil {
		if data.SWCalSub == "" {
			calsub = 0
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "版本号格式错误",
			})
			return
		}
	}

	if !dao.DBCheckBranchAndVersionRepeat(uint(branch), uint(version)) {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  fmt.Sprintf("没找到记录，检查软件分支号%d和版本号%d是否正确", branch, version),
		})
		return
	}

	var output string
	if data.SWBuildFile != nil {
		if data.SWA2LFile == nil || data.SWVersion == "" || data.SWModifier == "" || data.SWFinishTime == "" || data.SWLog == "" {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "SWA2LFile  SWVersion  SWModifier  SWFinishTime  SWLog 这些参数不能为空",
			})
			return
		}
		var epl model.ECUProjectList
		err := dao.DBWhereSWBranchFindECUProjectListRecord(uint(branch), &epl)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  err.Error(),
			})
			return
		}
		if data.SWBuildFile.Filename != epl.SoftwareName {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "软件名称不一致",
			})
			return
		}
		buildFilePath := "upload/" + data.SWBuildFile.Filename
		a2lFilePath := "upload/" + data.SWA2LFile.Filename

		_, err = os.Stat(buildFilePath)
		if err == nil {
			err = os.Remove(buildFilePath)
			if err != nil {
				c.JSON(400, serializer.Response{
					Code: 400,
					Msg:  "服务器删除文件时失败",
				})
				return
			}
		}

		err = c.SaveUploadedFile(data.SWBuildFile, buildFilePath)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "服务器写入文件时失败",
			})
			return
		} // 保存文件

		_, err = os.Stat(a2lFilePath)
		if err == nil {
			err = os.Remove(a2lFilePath)
			if err != nil {
				c.JSON(400, serializer.Response{
					Code: 400,
					Msg:  "服务器删除文件时失败",
				})
				return
			}
		}
		err = c.SaveUploadedFile(data.SWA2LFile, a2lFilePath)
		if err != nil {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "服务器写入文件时失败",
			})
			return
		}
		//}

		files := []string{buildFilePath, a2lFilePath}
		if data.SWDLLFile != nil {
			files = append(files, "upload/"+data.SWDLLFile.Filename)
		}
		output = "ECUSoftware/" + data.SWBuildFile.Filename[:len(data.SWBuildFile.Filename)-4] + ".zip"

		// 判断upload文件夹中是否存在指定的zip文件
		_, err = os.Stat(output)
		if err != nil {
			// 文件不存在，则将文件压缩打包
			if err = common.ZipFiles(output, files); err != nil {
				c.JSON(400, serializer.Response{
					Code: 400,
					Msg:  "压缩文件时失败",
				})
				return
			}
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  fmt.Sprintf("%s 文件已存在", output),
			})
			return
		}
	}

	changeAttached := ""
	sysVerifyAttached := ""
	cALFile := ""
	if data.ChangeAttached == nil {
		changeAttached = ""
	} else {
		changeAttached = data.ChangeAttached.Filename
	}
	if data.SysVerifyAttached == nil {
		sysVerifyAttached = ""
	} else {
		sysVerifyAttached = data.SysVerifyAttached.Filename
	}
	if data.CALFile == nil {
		cALFile = ""
	} else {
		cALFile = data.CALFile.Filename
	}
	if output != "" {
		output = output[12:]
	}
	ecusoftwareversion := model.ECUVer{
		ChangeInitiator:    data.ChangeInitiator,
		ChangeInitiateTime: data.ChangeInitiateTime,
		ChangeCause:        data.ChangeCause,
		ChangeReq:          data.ChangeReq,
		ChangeAttached:     changeAttached,
		ChangeApplyRange:   data.ChangeApplyRange,
		SWModifier:         data.SWModifier,
		SWFinishTime:       data.SWFinishTime,
		SWLog:              data.SWLog,
		SWBuildFile:        output,
		SWBranch:           uint(branch),
		SWVersion:          uint(version),
		SWCalMain:          uint(calmain),
		SWCalSub:           uint(calsub),
		HILTester:          data.HILTester,
		HILFinishTime:      data.HILFinishTime,
		HILResult:          data.HILResult,
		SysVerifier:        data.SysVerifier,
		SysVerifyTime:      data.SysVerifyTime,
		SysVerifyResult:    data.SysVerifyResult,
		SysVerifyAttached:  sysVerifyAttached,
		SWLevel:            4,
		SWLogClient:        data.SWLogClient,
		SWLogClientEN:      data.SWLogClientEN,
		CALFile:            cALFile,
		Comment:            data.Comment,
	}
	err = dao.ECUSoftwareVersionChange(adminId, ecusoftwareversion)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(200, serializer.Response{
		Code: 200,
		Msg:  "success",
	})
}

func VerRecordRelease(c *gin.Context) {
	var data ECUProjectRequestModel
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "参数绑定时出错",
		})
		return
	}
	adminIdAny, _ := c.Get("userId")
	adminId := adminIdAny.(uint)
	common.WriteLog(adminId, fmt.Sprintf("VerRecordRelease   软件分支号：%s  软件版本号：%s  软件等级：%s", data.SWBranch, data.SWVersion, data.SWLevel))
	branch, _ := strconv.ParseUint(data.SWBranch, 10, 64)
	version, err := strconv.ParseUint(data.SWVersion, 10, 64)
	if err != nil {
		if data.SWVersion == "" {
			version = 0
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "版本号格式错误",
			})
			return
		}
	}
	level, err := strconv.ParseUint(data.SWLevel, 10, 64)
	if err != nil {
		if data.SWLevel == "" {
			level = 0
		} else {
			c.JSON(400, serializer.Response{
				Code: 400,
				Msg:  "版本号格式错误",
			})
			return
		}
	}

	if !dao.DBCheckBranchAndVersionRepeat(uint(branch), uint(version)) {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  fmt.Sprintf("没找到记录，检查软件分支号%d和版本号%d是否正确", branch, version),
		})
		return
	}

	var ecuVer model.ECUVer
	err = dao.DBFindECUSoftwareRecord(uint(branch), uint(version), &ecuVer)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}
	if ecuVer.ChangeInitiator == "" || ecuVer.ChangeInitiateTime == "" || ecuVer.ChangeReq == "" || ecuVer.SWModifier == "" || ecuVer.SWFinishTime == "" || ecuVer.SWLog == "" || ecuVer.SWBuildFile == "" || ecuVer.SWLogClient == "" || ecuVer.SWLogClientEN == "" {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  "change_initiator\nchange_initiate_time\nchange_req\nsw_modifier\nsw_finish_time\nsw_log\nsw_build_file\nsw_log_client\nsw_log_client_en\n这些字段不能有空值",
		})
		return
	}

	e := model.ECUVer{
		SWBranch:  uint(branch),
		SWVersion: uint(version),
		SWLevel:   uint(level),
	}
	err = dao.ECUSoftwareVersionChange(adminId, e)
	if err != nil {
		c.JSON(400, serializer.Response{
			Code: 400,
			Msg:  err.Error(),
		})
		return
	}
}
