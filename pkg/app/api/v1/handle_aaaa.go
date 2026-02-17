package v1

import (
	"strconv"

	"github.com/cylonchau/hermes/pkg/app/api/query"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

type AAAARecordRouter struct {
	DAO *rdb.RecordDAO
}

func (ar *AAAARecordRouter) List(c *gin.Context) {
	var records []model.AAAARecord
	if err := model.DB.Preload("Record.Zone").Find(&records).Error; err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, records)
}

func (ar *AAAARecordRouter) Create(c *gin.Context) {
	var req struct {
		Record model.Record     `json:"record"`
		AAAA   model.AAAARecord `json:"aaaa"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		query.BadRequest(c, err)
		return
	}
	if err := ar.DAO.CreateAAAARecord(c.Request.Context(), &req.Record, &req.AAAA); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, req.AAAA)
}

func (ar *AAAARecordRouter) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := ar.DAO.GetAAAARecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (ar *AAAARecordRouter) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := ar.DAO.GetAAAARecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		query.BadRequest(c, err)
		return
	}

	if err := ar.DAO.UpdateAAAARecord(c.Request.Context(), &record.Record, record); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (ar *AAAARecordRouter) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if err := ar.DAO.DeleteAAAARecord(c.Request.Context(), uint(id)); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, nil)
}
