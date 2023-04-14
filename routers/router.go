package routers

import (
	"ConfigManager/controllers/request"
	"ConfigManager/globals"
)

func init() {

	api := globals.GvaGinEngine.Group("api")
	{
		config := api.Group("config")
		config.GET("list", request.ListConfig)
		config.POST("create", request.CreateConfig)
		config.POST("update", request.UpdateConfig)
		config.POST("delete", request.DeleteConfig)
		config.POST("history", request.HistoryConfig)
		config.POST("rollback", request.RollbackConfig)

		kube := api.Group("kube")
		kube.GET("namespace/list", request.ListNamespace)
	}

}
