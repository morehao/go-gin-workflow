package router

import (
	"go-gin-workflow/internal/workflow/controller/ctrProcess"

	"github.com/gin-gonic/gin"
)

// procInstRouter 初始化审批流程实例路由信息
func procInstRouter(routerGroup *gin.RouterGroup) {
	procInstCtr := ctrProcess.NewProcInstCtr()
	procInstGroup := routerGroup.Group("process")
	{
		procInstGroup.POST("start", procInstCtr.Start)      // 新建审批流程实例
		procInstGroup.POST("delete", procInstCtr.Delete)    // 删除审批流程实例
		procInstGroup.POST("update", procInstCtr.Update)    // 更新审批流程实例
		procInstGroup.GET("detail", procInstCtr.Detail)     // 根据ID获取审批流程实例
		procInstGroup.GET("pageList", procInstCtr.PageList) // 获取审批流程实例列表
	}
}
