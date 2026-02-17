package v1

import (
	"strconv"

	"github.com/cylonchau/hermes/pkg/app/api/query"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

type ARecordRouter struct {
	DAO *rdb.RecordDAO
}

func (ar *ARecordRouter) List(c *gin.Context) {
	// Note: Generic list for A records without filter is not directly in DAO GetARecords.
	// GetARecords expects zoneName and recordName.
	// For now, let's keep it simple or expand DAO later.
	var records []model.ARecord
	if err := model.DB.Preload("Record.Zone").Find(&records).Error; err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, records)
}

func (ar *ARecordRouter) Create(c *gin.Context) {
	var recordAReq struct {
		Record model.Record  `json:"record"`
		A      model.ARecord `json:"a"`
	}
	if err := c.ShouldBindJSON(&recordAReq); err != nil {
		query.BadRequest(c, err)
		return
	}
	if err := ar.DAO.CreateARecord(c.Request.Context(), &recordAReq.Record, &recordAReq.A); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, recordAReq.A)
}

func (ar *ARecordRouter) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := ar.DAO.GetARecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (ar *ARecordRouter) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := ar.DAO.GetARecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		query.BadRequest(c, err)
		return
	}

	if err := ar.DAO.UpdateARecord(c.Request.Context(), &record.Record, record); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (ar *ARecordRouter) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if err := ar.DAO.DeleteARecord(c.Request.Context(), uint(id)); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, nil)
}
