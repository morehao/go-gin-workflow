package objProcdef

import "go-gin-workflow/internal/workflow/object/objFlow"

type ProcdefBaseInfo struct {
	Name     string        `json:"name" form:"name" validate:"required"`         // 流程名称
	Resource *objFlow.Node `json:"resource" form:"resource" validate:"required"` // 流程配置
}
