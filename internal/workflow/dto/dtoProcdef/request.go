package dtoProcdef

import (
	"go-gin-workflow/internal/workflow/object/objCommon"
	"go-gin-workflow/internal/workflow/object/objProcdef"
)

type ProcdefSaveReq struct {
	objProcdef.ProcdefBaseInfo
}

type ProcdefDetailReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}

type ProcdefPageListReq struct {
	objCommon.PageQuery
}

type ProcdefDeleteReq struct {
	ID uint64 `json:"id" form:"id" validate:"required" label:"数据自增id"` // 数据自增id
}
