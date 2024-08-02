package errorCode

import "github.com/morehao/go-tools/gerror"

// 审批流程相关错误码，1002xx

var ProcdefSaveErr = gerror.Error{
	Code: 100200,
	Msg:  "保存审批流程失败",
}

var ProcdefDeleteErr = gerror.Error{
	Code: 100201,
	Msg:  "删除审批流程失败",
}

var ProcdefUpdateErr = gerror.Error{
	Code: 100202,
	Msg:  "修改审批流程失败",
}

var ProcdefGetDetailErr = gerror.Error{
	Code: 100203,
	Msg:  "查看审批流程失败",
}

var ProcdefGetPageListErr = gerror.Error{
	Code: 100204,
	Msg:  "查看审批流程列表失败",
}

var ProcdefNotExistErr = gerror.Error{
	Code: 100205,
	Msg:  "审批流程不存在",
}
