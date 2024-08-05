package errorCode

import "github.com/morehao/go-tools/gerror"

var ProcInstCreateErr = gerror.Error{
	Code: 100100,
	Msg:  "创建审批流程实例失败",
}

var ProcInstDeleteErr = gerror.Error{
	Code: 100101,
	Msg:  "删除审批流程实例失败",
}

var ProcInstUpdateErr = gerror.Error{
	Code: 100102,
	Msg:  "修改审批流程实例失败",
}

var ProcInstGetDetailErr = gerror.Error{
	Code: 100103,
	Msg:  "查看审批流程实例失败",
}

var ProcInstGetPageListErr = gerror.Error{
	Code: 100104,
	Msg:  "查看审批流程实例列表失败",
}

var ProcInstNotExistErr = gerror.Error{
	Code: 100105,
	Msg:  "审批流程实例不存在",
}
var CreatedPageListErr = gerror.Error{
	Code: 50000,
	Msg:  "我创建的流程实例分页列表失败",
}
var TodoPageListErr = gerror.Error{
	Code: 50000,
	Msg:  "待我审批的流程实例分页列表失败",
}
var NotifyPageListErr = gerror.Error{
	Code: 50000,
	Msg:  "抄送我的的流程实例分页列表失败",
}
