package request

import (
	"ConfigManager/controllers/response"
	"ConfigManager/models"
	"ConfigManager/services"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

var (
	cm services.ConfigMap
)

/*
	ListConfig
	获取所有/指定配置
	GET http://localhost:8009/api/config/list?configName=test1&pageSize=3&pageNum=1
*/
func ListConfig(c *gin.Context) {

	configName := c.Query("configName")
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	resp, err := cm.ListConfig(configName, pageSize, pageNum)
	if err != nil {
		response.Fail(c, err.Error(), "获取失败")
		return
	}
	response.Success(c, resp, "获取成功")

}

/*
	CreateConfig
	创建配置
	POST http://localhost:8009/api/config/create?isCreateCm=true
*/
func CreateConfig(c *gin.Context) {

	var req models.ConfigureStruct
	isCreateCm := c.Query("isCreateCm")
	_ = c.ShouldBindJSON(&req)
	var cm services.ConfigMap
	if err := cm.CreateConfig(req, isCreateCm); err != nil {
		response.Fail(c, err.Error(), "创建失败, 配置文件名称已存在")
		return
	}
	response.Success(c, nil, "创建成功")

}

/*
	UpdateConfig
	更新配置
	POST http://localhost:8009/api/config/update?configName=test13&isUpdateCm=true
*/
func UpdateConfig(c *gin.Context) {

	var (
		newReq models.ConfigureStruct
	)

	configName := c.Query("configName")
	isUpdateCm := c.Query("isUpdateCm")
	resp, _ := cm.GetConfig(configName)
	resp.Version = resp.ConfigName + "-" + time.Now().Format("200601020304")
	_ = c.ShouldBindJSON(&newReq)
	err := cm.UpdateConfig(&newReq, &resp, isUpdateCm)
	if err != nil {
		response.Fail(c, err.Error(), "配置更新失败")
		return
	}
	response.Success(c, "", "配置更新成功")

}

/*
	DeleteConfig
	删除配置
	POST http://localhost:8009/api/config/delete?configName=test9&configNamespace=test
*/
func DeleteConfig(c *gin.Context) {

	configName := c.Query("configName")
	configNamespace := c.Query("configNamespace")
	err := cm.DeleteConfig(configName, configNamespace)
	if err != nil {
		response.Fail(c, err.Error(), "删除失败, 配置文件名称不存在或已被删除")
		return
	}
	response.Success(c, nil, "删除成功")

}

/*
	HistoryConfig
	历史配置
	POST http://localhost:8009/api/config/history?configName=test9&pageSize=2&pageNum=1
*/
func HistoryConfig(c *gin.Context) {

	configName := c.Query("configName")
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	resp, err := cm.HistoryConfig(configName, pageSize, pageNum)
	if err != nil {
		response.Fail(c, err.Error(), "查询失败, 配置文件名称不存在或已被删除")
		return
	}
	response.Success(c, resp, "查询成功")

}

/*
	RollbackConfig
	回滚配置
	POST http://localhost:8009/api/config/rollback?configName=test13&version=test13-202304140625
*/
func RollbackConfig(c *gin.Context) {

	configName := c.Query("configName")
	version := c.Query("version")

	respOldcf, _ := cm.GetConfig(configName)
	respOldcf.Version = respOldcf.ConfigName + "-" + time.Now().Format("200601020304")

	info, err := cm.RollbackConfig(configName, version, &respOldcf)
	if err != nil {
		response.Fail(c, err.Error(), "回滚失败")
		return
	}
	response.Success(c, info, "回滚成功")

}
