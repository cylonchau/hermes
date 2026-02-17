package v1

import (
	"strconv"

	"github.com/cylonchau/hermes/pkg/app/api/query"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

type SRVRecordRouter struct {
	DAO *rdb.RecordDAO
}

func (sr *SRVRecordRouter) List(c *gin.Context) {
	var records []model.SRVRecord
	if err := model.DB.Preload("Record.Zone").Find(&records).Error; err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, records)
}

func (sr *SRVRecordRouter) Create(c *gin.Context) {
	var req struct {
		Record model.Record    `json:"record"`
		SRV    model.SRVRecord `json:"srv"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		query.BadRequest(c, err)
		return
	}
	if err := sr.DAO.CreateSRVRecord(c.Request.Context(), &req.Record, &req.SRV); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, req.SRV)
}

func (sr *SRVRecordRouter) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := sr.DAO.GetSRVRecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (sr *SRVRecordRouter) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := sr.DAO.GetSRVRecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		query.BadRequest(c, err)
		return
	}

	if err := sr.DAO.UpdateSRVRecord(c.Request.Context(), &record.Record, record); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (sr *SRVRecordRouter) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if err := sr.DAO.DeleteSRVRecord(c.Request.Context(), uint(id)); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, nil)
}
