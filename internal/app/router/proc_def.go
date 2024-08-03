package router

import (
	"go-gin-workflow/internal/app/controller/ctrProcDef"

	"github.com/gin-gonic/gin"
)

// procDefRouter 初始化审批流程定义路由信息
func procDefRouter(routerGroup *gin.RouterGroup) {
	procDefCtr := ctrProcDef.NewProcDefCtr()
	procDefGroup := routerGroup.Group("procDef")
	{
		procDefGroup.POST("create", procDefCtr.Save)      // 新建审批流程定义
		procDefGroup.POST("delete", procDefCtr.Delete)    // 删除审批流程定义
		procDefGroup.GET("detail", procDefCtr.Detail)     // 根据ID获取审批流程定义
		procDefGroup.GET("pageList", procDefCtr.PageList) // 获取审批流程定义列表
	}
}
