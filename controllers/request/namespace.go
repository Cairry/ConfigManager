package request

import (
	"ConfigManager/controllers/response"
	"ConfigManager/services"
	"github.com/gin-gonic/gin"
)

func ListNamespace(c *gin.Context) {

	info, err := services.ListNamespace()
	if err != nil {
		response.Fail(c, err.Error(), "获取失败")
		return
	}
	response.Success(c, info, "获取成功")

}
