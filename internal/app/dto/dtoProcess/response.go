package dtoProcess

import (
	"go-gin-workflow/internal/app/object/objCommon"
	"go-gin-workflow/internal/app/object/objProcess"
)

type ProcInstStartResp struct {
	ID uint64 `json:"id"` // 数据自增id
}

type ProcInstDetailResp struct {
	ID uint64 `json:"id" validate:"required"` // 数据自增id
	objProcess.ProcInstBaseInfo
	objCommon.OperatorBaseInfo
}

type ProcInstPageListItem struct {
	ID uint64 `json:"id" validate:"required"` // 数据自增id
	objProcess.ProcInstBaseInfo
	objCommon.OperatorBaseInfo
}

type ProcInstPageListResp struct {
	List  []ProcInstPageListItem `json:"list"`  // 数据列表
	Total int64                  `json:"total"` // 数据总条数
}
