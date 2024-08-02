package daoProcdef

import (
	"fmt"
	"go-gin-workflow/internal/pkg/errorCode"
	"go-gin-workflow/internal/workflow/model"
	"go-gin-workflow/internal/workflow/object/objFlow"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
)

// ProcdefHistoryEntity 审批流程定义历史记录表结构体
type ProcdefHistoryEntity struct {
	ID         uint64        `gorm:"column:id;comment:自增id;primaryKey"`
	Company    string        `gorm:"column:company;comment:公司名称"`
	Name       string        `gorm:"column:name;comment:流程名称"`
	Resource   *objFlow.Node `gorm:"column:resource;type:json;comment:流程配置"`
	UserID     string        `gorm:"column:userid;comment:用户id"`
	Username   string        `gorm:"column:username;comment:用户名称"`
	DeployTime string        `gorm:"column:deploy_time;comment:部署时间"`
	Version    uint64        `gorm:"column:version;comment:流程版本"`
}

type ProcdefHistoryEntityList []ProcdefHistoryEntity

const TblNameProcdefHistory = "procdef_history"

func (ProcdefHistoryEntity) TableName() string {
	return TblNameProcdefHistory
}

type ProcdefHistoryCond struct {
	ID             uint64
	IDs            []uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type ProcdefHistoryDao struct {
	model.Base
}

func NewProcdefHistoryDao() *ProcdefHistoryDao {
	return &ProcdefHistoryDao{}
}

func (dao *ProcdefHistoryDao) WithTx(db *gorm.DB) *ProcdefHistoryDao {
	dao.Tx = db
	return dao
}

func (dao *ProcdefHistoryDao) Insert(c *gin.Context, entity *ProcdefHistoryEntity) error {
	db := dao.Db(c).Model(&ProcdefHistoryEntity{})
	db = db.Table(TblNameProcdefHistory)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ProcdefHistoryDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ProcdefHistoryDao) BatchInsert(c *gin.Context, entityList ProcdefHistoryEntityList) error {
	db := dao.Db(c).Table(TblNameProcdefHistory)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ProcdefHistoryDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *ProcdefHistoryDao) Update(c *gin.Context, entity *ProcdefHistoryEntity) error {
	db := dao.Db(c).Model(&ProcdefHistoryEntity{})
	db = db.Table(TblNameProcdefHistory)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcdefHistoryDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ProcdefHistoryDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&ProcdefHistoryEntity{})
	db = db.Table(TblNameProcdefHistory)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcdefHistoryDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *ProcdefHistoryDao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&ProcdefHistoryEntity{})
	db = db.Table(TblNameProcdefHistory)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcdefHistoryDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *ProcdefHistoryDao) GetById(c *gin.Context, id uint64) (*ProcdefHistoryEntity, error) {
	var entity ProcdefHistoryEntity
	db := dao.Db(c).Model(&ProcdefHistoryEntity{})
	db = db.Table(TblNameProcdefHistory)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcdefHistoryDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *ProcdefHistoryDao) GetByCond(c *gin.Context, cond *ProcdefHistoryCond) (*ProcdefHistoryEntity, error) {
	var entity ProcdefHistoryEntity
	db := dao.Db(c).Model(&ProcdefHistoryEntity{})
	db = db.Table(TblNameProcdefHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcdefHistoryDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *ProcdefHistoryDao) GetListByCond(c *gin.Context, cond *ProcdefHistoryCond) (ProcdefHistoryEntityList, error) {
	var entityList ProcdefHistoryEntityList
	db := dao.Db(c).Model(&ProcdefHistoryEntity{})
	db = db.Table(TblNameProcdefHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcdefHistoryDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *ProcdefHistoryDao) GetPageListByCond(c *gin.Context, cond *ProcdefHistoryCond) (ProcdefHistoryEntityList, int64, error) {
	db := dao.Db(c).Model(&ProcdefHistoryEntity{})
	db = db.Table(TblNameProcdefHistory)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ProcdefHistoryDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list ProcdefHistoryEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ProcdefHistoryDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l ProcdefHistoryEntityList) ToMap() map[uint64]ProcdefHistoryEntity {
	m := make(map[uint64]ProcdefHistoryEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *ProcdefHistoryDao) BuildCondition(db *gorm.DB, cond *ProcdefHistoryCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameProcdefHistory)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameProcdefHistory)
		db.Where(query, cond.IDs)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameProcdefHistory)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameProcdefHistory)
		db.Where(query, time.Unix(cond.CreatedAtEnd, 0))
	}
	if cond.IsDelete {
		db.Unscoped()
	}

	if cond.OrderField != "" {
		db.Order(cond.OrderField)
	}

	return
}
