package router

import (
	"go-gin-workflow/internal/workflow/controller/ctrProcdef"

	"github.com/gin-gonic/gin"
)

// procdefRouter 初始化审批流程定义路由信息
func procdefRouter(routerGroup *gin.RouterGroup) {
	procdefCtr := ctrProcdef.NewProcdefCtr()
	procdefGroup := routerGroup.Group("procdef")
	{
		procdefGroup.POST("create", procdefCtr.Save)      // 新建审批流程定义
		procdefGroup.POST("delete", procdefCtr.Delete)    // 删除审批流程定义
		procdefGroup.GET("detail", procdefCtr.Detail)     // 根据ID获取审批流程定义
		procdefGroup.GET("pageList", procdefCtr.PageList) // 获取审批流程定义列表
	}
}
