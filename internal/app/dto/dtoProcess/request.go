package dtoProcess

import (
	"go-gin-workflow/internal/app/object/objCommon"
	"go-gin-workflow/internal/app/object/objProcess"
)

type ProcInstStartReq struct {
	ProcDefName string            `json:"procDefName" validate:"required"` // 流程定义名称
	Title       string            `json:"title" validate:"required"`       // 流程标题
	Var         map[string]string `json:"var"`                             // 流程变量
}

type ProcInstUpdateReq struct {
	ID uint64 `json:"id" validate:"required" label:"数据自增id"` // 数据自增id
	objProcess.ProcInstBaseInfo
}

type ProcInstDetailReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type ProcInstPageListReq struct {
	objCommon.PageQuery
}

type ProcInstDeleteReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}
type CreatedPageListReq struct {
	objCommon.PageQuery
}
type TodoPageListReq struct {
	objCommon.PageQuery
	GroupList      []string `json:"groupList" form:"groupList"`           // 用户组列表
	DepartmentList []string `json:"departmentList" form:"departmentList"` // 部门列表
}
type NotifyPageListReq struct {
	objCommon.PageQuery
	GroupList []string `json:"groupList" form:"groupList"` // 用户组列表
}
