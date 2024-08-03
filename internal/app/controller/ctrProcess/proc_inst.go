package ctrProcess

import (
	"go-gin-workflow/internal/app/dto/dtoProcess"
	"go-gin-workflow/internal/app/service/svcProcess"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gcontext/ginRender"
)

type ProcInstCtr interface {
	Start(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
	Detail(c *gin.Context)
	PageList(c *gin.Context)
}

type procInstCtr struct {
	procInstSvc svcProcess.ProcInstSvc
}

var _ ProcInstCtr = (*procInstCtr)(nil)

func NewProcInstCtr() ProcInstCtr {
	return &procInstCtr{
		procInstSvc: svcProcess.NewProcInstSvc(),
	}
}

// Start 启动审批流程实例
// @Tags 审批流程实例
// @Summary 启动审批流程实例
// @accept application/json
// @Produce application/json
// @Param req body dtoProcess.ProcInstStartReq true "启动审批流程实例"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.ProcInstStartResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/process/start [post]
func (ctr *procInstCtr) Start(c *gin.Context) {
	var req dtoProcess.ProcInstStartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.Start(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// Delete 删除审批流程实例
// @Tags 审批流程实例
// @Summary 删除审批流程实例
// @accept application/json
// @Produce application/json
// @Param req body dtoProcess.ProcInstDeleteReq true "删除审批流程实例"
// @Success 200 {object} dto.DefaultRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /workflow/process/delete [post]
func (ctr *procInstCtr) Delete(c *gin.Context) {
	var req dtoProcess.ProcInstDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}

	if err := ctr.procInstSvc.Delete(c, &req); err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, "删除成功")
	}
}

// Update 修改审批流程实例
// @Tags 审批流程实例
// @Summary 修改审批流程实例
// @accept application/json
// @Produce application/json
// @Param req body dtoProcess.ProcInstUpdateReq true "修改审批流程实例"
// @Success 200 {object} dto.DefaultRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /workflow/process/update [post]
func (ctr *procInstCtr) Update(c *gin.Context) {
	var req dtoProcess.ProcInstUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	if err := ctr.procInstSvc.Update(c, &req); err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, "修改成功")
	}
}

// Detail 审批流程实例详情
// @Tags 审批流程实例
// @Summary 审批流程实例详情
// @accept application/json
// @Produce application/json
// @Param req query dtoProcess.ProcInstDetailReq true "审批流程实例详情"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.ProcInstDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/process/detail [get]
func (ctr *procInstCtr) Detail(c *gin.Context) {
	var req dtoProcess.ProcInstDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.Detail(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// PageList 审批流程实例列表
// @Tags 审批流程实例
// @Summary 审批流程实例列表分页
// @accept application/json
// @Produce application/json
// @Param req query dtoProcess.ProcInstPageListReq true "审批流程实例列表"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.ProcInstPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/process/pageList [get]
func (ctr *procInstCtr) PageList(c *gin.Context) {
	var req dtoProcess.ProcInstPageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.PageList(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}
