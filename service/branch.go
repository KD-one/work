package service

import (
	"github.com/gin-gonic/gin"
	"test/common"
	"test/model"
)

func BranchToVersion1(c *gin.Context) {
	c.HTML(200, "branchToVersion/index.html", nil)
}

func BranchToVersion2(c *gin.Context) {
	branch := c.Query("SoftwareVersionBranch")
	var f []model.Filemap
	common.DB.Model(model.Filemap{}).Where("software_version_branch = ?", branch).Order("software_version_number").Find(&f)

	// versions 根据版本号分组   版本号:文件名数组
	versions := make(map[string][]string)
	var fileName string
	for _, software := range f {
		fileName = branch + "_" + software.SoftwareVersionNumber + "_" + software.SoftwareName + ".zip"
		versions[software.SoftwareVersionNumber] = append(versions[software.SoftwareVersionNumber], fileName)

	}

	c.HTML(200, "branchToVersion/btv.html", gin.H{
		"branch":   branch,
		"versions": versions,
	})
}
