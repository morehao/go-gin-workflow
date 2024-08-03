package svcProcess

import (
	"go-gin-workflow/internal/pkg/constants"
	"go-gin-workflow/internal/pkg/context"
	"go-gin-workflow/internal/pkg/errorCode"
	"go-gin-workflow/internal/workflow/dto/dtoProcess"
	"go-gin-workflow/internal/workflow/flow"
	"go-gin-workflow/internal/workflow/helper"
	"go-gin-workflow/internal/workflow/model/daoProcdef"
	"go-gin-workflow/internal/workflow/model/daoProcess"
	"go-gin-workflow/internal/workflow/object/objCommon"
	"go-gin-workflow/internal/workflow/object/objFlow"
	"go-gin-workflow/internal/workflow/object/objProcess"

	jsoniter "github.com/json-iterator/go"

	"gorm.io/gorm"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
)

type ProcInstSvc interface {
	Start(c *gin.Context, req *dtoProcess.ProcInstStartReq) (*dtoProcess.ProcInstStartResp, error)
	Delete(c *gin.Context, req *dtoProcess.ProcInstDeleteReq) error
	Update(c *gin.Context, req *dtoProcess.ProcInstUpdateReq) error
	Detail(c *gin.Context, req *dtoProcess.ProcInstDetailReq) (*dtoProcess.ProcInstDetailResp, error)
	PageList(c *gin.Context, req *dtoProcess.ProcInstPageListReq) (*dtoProcess.ProcInstPageListResp, error)
}

type procInstSvc struct {
}

var _ ProcInstSvc = (*procInstSvc)(nil)

func NewProcInstSvc() ProcInstSvc {
	return &procInstSvc{}
}

// Start 创建审批流程实例
func (svc *procInstSvc) Start(c *gin.Context, req *dtoProcess.ProcInstStartReq) (*dtoProcess.ProcInstStartResp, error) {
	companyName, userID, userName, departmentName := context.GetCompanyName(c), context.GetUserID(c), context.GetUserName(c), context.GetDepartmentName(c)

	procDefEntity, getProcDefErr := daoProcdef.NewProcdefDao().GetByCond(c, &daoProcdef.ProcdefCond{
		Company: companyName,
		Name:    req.ProcDefName,
	})
	if getProcDefErr != nil {
		glog.Errorf(c, "[svcProcess.ProcInstCreate] daoProcdef GetByCond fail, err:%v, req:%s", getProcDefErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}
	if procDefEntity == nil || procDefEntity.ID == 0 {
		return nil, errorCode.ProcdefNotExistErr
	}
	now := time.Now()
	instEntity := &daoProcess.ProcInstEntity{
		ProcDefID:     procDefEntity.ID,
		ProcDefName:   procDefEntity.Name,
		NodeID:        procDefEntity.Resource.NodeID,
		Title:         req.Title,
		Company:       companyName,
		Department:    departmentName,
		StartUserID:   userID,
		StartUserName: userName,
		StartTime:     gutils.TimeFormat(now, gutils.YYYY_MM_DD_HH_MM_SS),
	}
	execNodeLinkedList, parseErr := flow.ParseProcessConfig(procDefEntity.Resource, req.Var)
	if parseErr != nil {
		glog.Errorf(c, "[svcProcess.ProcInstCreate] parse process config fail, err:%v, req:%s", parseErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}
	execNodeLinkedList.PushBack(objFlow.ExecNode{
		NodeID: constants.NodeIDTextEnd,
	})
	execNodeLinkedList.PushFront(objFlow.ExecNode{
		NodeID:   constants.NodeIDTextStart,
		Type:     constants.NodeInfoTypeStarter,
		Approver: userID,
	})
	execNodeList := gutils.LinkedListToArray(execNodeLinkedList)

	executionEntity := &daoProcess.ExecutionEntity{
		ProcDefID: procDefEntity.ID,
		NodeInfos: gutils.ToJsonString(execNodeList),
	}
	var execNodeStList []objFlow.ExecNode
	if err := jsoniter.Unmarshal([]byte(executionEntity.NodeInfos), &execNodeStList); err != nil {
		glog.Errorf(c, "[svcProcess.ProcInstCreate] unmarshal execution node info fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}

	taskEntity := &daoProcess.TaskEntity{
		NodeID:        constants.NodeIDTextStart,
		Assignee:      userID,
		IsFinished:    constants.ProcessTaskStatusFinished,
		ClaimTime:     gutils.TimeFormat(now, gutils.YYYY_MM_DD_HH_MM_SS),
		Step:          0,
		MemberCount:   1,
		UnCompleteNum: 0,
		ActType:       constants.ActionTypeOr,
		AgreeNum:      1,
	}
	if len(execNodeStList) > 0 {
		if execNodeStList[0].ActType == constants.ActionTypeAnd {
			taskEntity.UnCompleteNum = execNodeStList[0].MemberCount
			taskEntity.MemberCount = execNodeStList[0].MemberCount
		}
	}

	txErr := helper.MysqlClient.Transaction(func(tx *gorm.DB) error {
		if err := daoProcess.NewProcInstDao().WithTx(tx).Insert(c, instEntity); err != nil {
			return err
		}
		executionEntity.ProcInstID = instEntity.ID
		if err := daoProcess.NewExecutionDao().Insert(c, executionEntity); err != nil {
			return err
		}
		taskEntity.ProcInstID = instEntity.ID
		if err := daoProcess.NewTaskDao().WithTx(tx).Insert(c, taskEntity); err != nil {
			return err
		}
		moveStageParam := &flow.MoveStageParam{
			ExecNodeList: execNodeStList,
			UserID:       userID,
			Company:      companyName,
			ProcInstID:   instEntity.ID,
			Comment:      "开始流程",
			TaskID:       taskEntity.ID,
			Step:         0,
			Pass:         true,
		}
		if err := flow.MoveStage(c, tx, moveStageParam); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(c, "[svcProcess.ProcInstCreate] tx fail, err:%v, req:%s", txErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}
	if err := daoProcess.NewProcInstDao().Insert(c, instEntity); err != nil {
		glog.Errorf(c, "[svcProcess.ProcInstCreate] daoProcInst Start fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}

	return &dtoProcess.ProcInstStartResp{
		ID: instEntity.ID,
	}, nil
}

// Delete 删除审批流程实例
func (svc *procInstSvc) Delete(c *gin.Context, req *dtoProcess.ProcInstDeleteReq) error {
	if err := daoProcess.NewProcInstDao().Delete(c, req.ID, 0); err != nil {
		glog.Errorf(c, "[svcProcess.Delete] daoProcInst Delete fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return errorCode.ProcInstDeleteErr
	}
	return nil
}

// Update 更新审批流程实例
func (svc *procInstSvc) Update(c *gin.Context, req *dtoProcess.ProcInstUpdateReq) error {
	updateEntity := &daoProcess.ProcInstEntity{
		ID: req.ID,
	}
	if err := daoProcess.NewProcInstDao().Update(c, updateEntity); err != nil {
		glog.Errorf(c, "[svcProcess.ProcInstUpdate] daoProcInst Update fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return errorCode.ProcInstUpdateErr
	}
	return nil
}

// Detail 根据id获取审批流程实例
func (svc *procInstSvc) Detail(c *gin.Context, req *dtoProcess.ProcInstDetailReq) (*dtoProcess.ProcInstDetailResp, error) {
	detailEntity, err := daoProcess.NewProcInstDao().GetById(c, req.ID)
	if err != nil {
		glog.Errorf(c, "[svcProcess.ProcInstDetail] daoProcInst GetById fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstGetDetailErr
	}
	// 判断是否存在
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, errorCode.ProcInstNotExistErr
	}
	Resp := &dtoProcess.ProcInstDetailResp{
		ID: detailEntity.ID,
		ProcInstBaseInfo: objProcess.ProcInstBaseInfo{
			Candidate:     detailEntity.Candidate,
			Company:       detailEntity.Company,
			Department:    detailEntity.Department,
			Duration:      detailEntity.Duration,
			EndTime:       detailEntity.EndTime,
			IsFinished:    detailEntity.IsFinished,
			NodeID:        detailEntity.NodeID,
			ProcDefID:     detailEntity.ProcDefID,
			ProcDefName:   detailEntity.ProcDefName,
			StartTime:     detailEntity.StartTime,
			StartUserID:   detailEntity.StartUserID,
			StartUserName: detailEntity.StartUserName,
			TaskID:        detailEntity.TaskID,
			Title:         detailEntity.Title,
		},
		OperatorBaseInfo: objCommon.OperatorBaseInfo{},
	}
	return Resp, nil
}

// PageList 分页获取审批流程实例列表
func (svc *procInstSvc) PageList(c *gin.Context, req *dtoProcess.ProcInstPageListReq) (*dtoProcess.ProcInstPageListResp, error) {
	cond := &daoProcess.ProcInstCond{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	dataList, total, err := daoProcess.NewProcInstDao().GetPageListByCond(c, cond)
	if err != nil {
		glog.Errorf(c, "[svcProcess.ProcInstPageList] daoProcInst GetPageListByCond fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstGetPageListErr
	}
	list := make([]dtoProcess.ProcInstPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtoProcess.ProcInstPageListItem{
			ID: v.ID,
			ProcInstBaseInfo: objProcess.ProcInstBaseInfo{
				Candidate:     v.Candidate,
				Company:       v.Company,
				Department:    v.Department,
				Duration:      v.Duration,
				EndTime:       v.EndTime,
				IsFinished:    v.IsFinished,
				NodeID:        v.NodeID,
				ProcDefID:     v.ProcDefID,
				ProcDefName:   v.ProcDefName,
				StartTime:     v.StartTime,
				StartUserID:   v.StartUserID,
				StartUserName: v.StartUserName,
				TaskID:        v.TaskID,
				Title:         v.Title,
			},
			OperatorBaseInfo: objCommon.OperatorBaseInfo{},
		})
	}
	return &dtoProcess.ProcInstPageListResp{
		List:  list,
		Total: total,
	}, nil
}
