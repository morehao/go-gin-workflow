package ctrProcdef

import (
	"go-gin-workflow/internal/workflow/dto/dtoProcdef"
	"go-gin-workflow/internal/workflow/service/svcProcdef"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gcontext/ginRender"
)

type ProcdefCtr interface {
	Save(c *gin.Context)
	Delete(c *gin.Context)
	Detail(c *gin.Context)
	PageList(c *gin.Context)
}

type procdefCtr struct {
	procdefSvc svcProcdef.ProcdefSvc
}

var _ ProcdefCtr = (*procdefCtr)(nil)

func NewProcdefCtr() ProcdefCtr {
	return &procdefCtr{
		procdefSvc: svcProcdef.NewProcdefSvc(),
	}
}

// Save 创建审批流程定义
// @Tags 审批流程定义
// @Summary 创建审批流程定义
// @accept application/json
// @Produce application/json
// @Param req body dtoProcdef.ProcdefSaveReq true "创建审批流程定义"
// @Success 200 {object} dto.DefaultRender{data=dtoProcdef.ProcdefSaveResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procdef/save [post]
func (ctr *procdefCtr) Save(c *gin.Context) {
	var req dtoProcdef.ProcdefSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procdefSvc.Save(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// Delete 删除审批流程定义
// @Tags 审批流程定义
// @Summary 删除审批流程定义
// @accept application/json
// @Produce application/json
// @Param req body dtoProcdef.ProcdefDeleteReq true "删除审批流程定义"
// @Success 200 {object} dto.DefaultRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /workflow/procdef/delete [post]
func (ctr *procdefCtr) Delete(c *gin.Context) {
	var req dtoProcdef.ProcdefDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}

	if err := ctr.procdefSvc.Delete(c, &req); err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, "删除成功")
	}
}

// Detail 审批流程定义详情
// @Tags 审批流程定义
// @Summary 审批流程定义详情
// @accept application/json
// @Produce application/json
// @Param req query dtoProcdef.ProcdefDetailReq true "审批流程定义详情"
// @Success 200 {object} dto.DefaultRender{data=dtoProcdef.ProcdefDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procdef/detail [get]
func (ctr *procdefCtr) Detail(c *gin.Context) {
	var req dtoProcdef.ProcdefDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procdefSvc.Detail(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.SuccessWithFormat(c, res)
	}
}

// PageList 审批流程定义列表
// @Tags 审批流程定义
// @Summary 审批流程定义列表分页
// @accept application/json
// @Produce application/json
// @Param req query dtoProcdef.ProcdefPageListReq true "审批流程定义列表"
// @Success 200 {object} dto.DefaultRender{data=dtoProcdef.ProcdefPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procdef/pageList [get]
func (ctr *procdefCtr) PageList(c *gin.Context) {
	var req dtoProcdef.ProcdefPageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procdefSvc.PageList(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}
