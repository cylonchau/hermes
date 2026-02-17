package v1

import (
	"strconv"

	"github.com/cylonchau/hermes/pkg/app/api/query"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

type MXRecordRouter struct {
	DAO *rdb.RecordDAO
}

func (mr *MXRecordRouter) List(c *gin.Context) {
	var records []model.MXRecord
	if err := model.DB.Preload("Record.Zone").Find(&records).Error; err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, records)
}

func (mr *MXRecordRouter) Create(c *gin.Context) {
	var req struct {
		Record model.Record   `json:"record"`
		MX     model.MXRecord `json:"mx"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		query.BadRequest(c, err)
		return
	}
	if err := mr.DAO.CreateMXRecord(c.Request.Context(), &req.Record, &req.MX); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, req.MX)
}

func (mr *MXRecordRouter) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := mr.DAO.GetMXRecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (mr *MXRecordRouter) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := mr.DAO.GetMXRecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		query.BadRequest(c, err)
		return
	}

	if err := mr.DAO.UpdateMXRecord(c.Request.Context(), &record.Record, record); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (mr *MXRecordRouter) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if err := mr.DAO.DeleteMXRecord(c.Request.Context(), uint(id)); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, nil)
}
