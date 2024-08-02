package svcProcdef

import (
	"go-gin-workflow/internal/pkg/context"
	"go-gin-workflow/internal/pkg/errorCode"
	"go-gin-workflow/internal/workflow/dto/dtoProcdef"
	"go-gin-workflow/internal/workflow/flow"
	"go-gin-workflow/internal/workflow/helper"
	"go-gin-workflow/internal/workflow/model/daoProcdef"
	"go-gin-workflow/internal/workflow/object/objCommon"
	"go-gin-workflow/internal/workflow/object/objProcdef"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
)

type ProcdefSvc interface {
	Save(c *gin.Context, req *dtoProcdef.ProcdefSaveReq) (*dtoProcdef.ProcdefSaveResp, error)
	Delete(c *gin.Context, req *dtoProcdef.ProcdefDeleteReq) error
	Detail(c *gin.Context, req *dtoProcdef.ProcdefDetailReq) (*dtoProcdef.ProcdefDetailResp, error)
	PageList(c *gin.Context, req *dtoProcdef.ProcdefPageListReq) (*dtoProcdef.ProcdefPageListResp, error)
}

type procdefSvc struct {
}

var _ ProcdefSvc = (*procdefSvc)(nil)

func NewProcdefSvc() ProcdefSvc {
	return &procdefSvc{}
}

// Save 创建审批流程定义
func (svc *procdefSvc) Save(c *gin.Context, req *dtoProcdef.ProcdefSaveReq) (*dtoProcdef.ProcdefSaveResp, error) {
	if err := flow.IsValidProcessConfig(req.Resource); err != nil {
		glog.Errorf(c, "[svcProcdef.Save] flow.IsValidProcessConfig fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcdefSaveErr
	}

	companyName, userId, userName := context.GetCompanyName(c), context.GetDepartmentName(c), context.GetUserName(c)

	maxVersionEntity, getMaxVersionErr := daoProcdef.NewProcdefDao().GetByCond(c, &daoProcdef.ProcdefCond{
		Company:    companyName,
		Name:       req.Name,
		OrderField: "id desc",
	})
	if getMaxVersionErr != nil {
		glog.Errorf(c, "[svcProcdef.Save] daoProcdef GetByCond fail, err:%v, req:%s", getMaxVersionErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcdefSaveErr
	}
	var version uint64 = 1
	if maxVersionEntity != nil && maxVersionEntity.ID > 0 {
		version = maxVersionEntity.Version + 1
	}
	now := time.Now()
	insertEntity := &daoProcdef.ProcdefEntity{
		Company:    companyName,
		Name:       req.Name,
		Resource:   req.Resource,
		UserID:     userId,
		Username:   userName,
		DeployTime: gutils.TimeFormat(now, gutils.YYYY_MM_DD_HH_MM_SS),
		Version:    version,
	}
	txErr := helper.MysqlClient.Transaction(func(tx *gorm.DB) error {
		// 创建审批流程定义
		if err := daoProcdef.NewProcdefDao().WithTx(tx).Insert(c, insertEntity); err != nil {
			return err
		}
		if maxVersionEntity == nil || maxVersionEntity.ID <= 0 {
			return nil
		}
		// 迁移历史版本&&删除旧版本
		historyInsertEntity := &daoProcdef.ProcdefHistoryEntity{
			ID:         maxVersionEntity.ID,
			Company:    maxVersionEntity.Company,
			Name:       maxVersionEntity.Name,
			Resource:   maxVersionEntity.Resource,
			UserID:     maxVersionEntity.UserID,
			Username:   maxVersionEntity.Username,
			DeployTime: maxVersionEntity.DeployTime,
			Version:    maxVersionEntity.Version,
		}
		if err := daoProcdef.NewProcdefHistoryDao().WithTx(tx).Insert(c, historyInsertEntity); err != nil {
			return err
		}
		if err := daoProcdef.NewProcdefDao().WithTx(tx).Delete(c, maxVersionEntity.ID); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(c, "[svcProcdef.Save] txErr:%v, req:%s", txErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcdefSaveErr
	}
	return &dtoProcdef.ProcdefSaveResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除审批流程定义
func (svc *procdefSvc) Delete(c *gin.Context, req *dtoProcdef.ProcdefDeleteReq) error {

	defEntity, getEntityErr := daoProcdef.NewProcdefDao().GetById(c, req.ID)
	if getEntityErr != nil {
		glog.Errorf(c, "[svcProcdef.Delete] daoProcdef GetById fail, err:%v, req:%s", getEntityErr, gutils.ToJsonString(req))
		return errorCode.ProcdefGetDetailErr
	}
	if defEntity == nil || defEntity.ID <= 0 {
		return errorCode.ProcdefNotExistErr
	}

	if err := daoProcdef.NewProcdefDao().Delete(c, req.ID); err != nil {
		glog.Errorf(c, "[svcProcdef.Delete] daoProcdef Delete fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return errorCode.ProcdefDeleteErr
	}
	return nil
}

// Detail 根据id获取审批流程定义
func (svc *procdefSvc) Detail(c *gin.Context, req *dtoProcdef.ProcdefDetailReq) (*dtoProcdef.ProcdefDetailResp, error) {
	detailEntity, err := daoProcdef.NewProcdefDao().GetById(c, req.ID)
	if err != nil {
		glog.Errorf(c, "[svcProcdef.ProcdefDetail] daoProcdef GetById fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcdefGetDetailErr
	}
	// 判断是否存在
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, errorCode.ProcdefNotExistErr
	}
	Resp := &dtoProcdef.ProcdefDetailResp{
		ID: detailEntity.ID,
		ProcdefBaseInfo: objProcdef.ProcdefBaseInfo{
			Name:     detailEntity.Name,
			Resource: detailEntity.Resource,
		},
		DeployTime: detailEntity.DeployTime,
		Version:    detailEntity.Version,
		OperatorBaseInfo: objCommon.OperatorBaseInfo{
			Company:  detailEntity.Company,
			UserID:   detailEntity.UserID,
			UserName: detailEntity.Username,
		},
	}
	return Resp, nil
}

// PageList 分页获取审批流程定义列表
func (svc *procdefSvc) PageList(c *gin.Context, req *dtoProcdef.ProcdefPageListReq) (*dtoProcdef.ProcdefPageListResp, error) {
	cond := &daoProcdef.ProcdefCond{
		Page:       req.Page,
		PageSize:   req.PageSize,
		OrderField: "id desc",
	}
	dataList, total, err := daoProcdef.NewProcdefDao().GetPageListByCond(c, cond)
	if err != nil {
		glog.Errorf(c, "[svcProcdef.ProcdefPageList] daoProcdef GetPageListByCond fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcdefGetPageListErr
	}
	list := make([]dtoProcdef.ProcdefPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtoProcdef.ProcdefPageListItem{
			ID:         v.ID,
			Name:       v.Name,
			Company:    v.Company,
			DeployTime: v.DeployTime,
			Username:   v.Username,
			UserID:     v.UserID,
		})
	}
	return &dtoProcdef.ProcdefPageListResp{
		List:  list,
		Total: total,
	}, nil
}
