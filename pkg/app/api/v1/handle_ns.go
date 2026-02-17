package v1

import (
	"strconv"

	"github.com/cylonchau/hermes/pkg/app/api/query"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

type NSRecordRouter struct {
	DAO *rdb.RecordDAO
}

func (nr *NSRecordRouter) List(c *gin.Context) {
	var records []model.NSRecord
	if err := model.DB.Preload("Record.Zone").Find(&records).Error; err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, records)
}

func (nr *NSRecordRouter) Create(c *gin.Context) {
	var req struct {
		Record model.Record   `json:"record"`
		NS     model.NSRecord `json:"ns"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		query.BadRequest(c, err)
		return
	}
	if err := nr.DAO.CreateNSRecord(c.Request.Context(), &req.Record, &req.NS); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, req.NS)
}

func (nr *NSRecordRouter) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := nr.DAO.GetNSRecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (nr *NSRecordRouter) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	record, err := nr.DAO.GetNSRecordByID(c.Request.Context(), uint(id))
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		query.BadRequest(c, err)
		return
	}

	if err := nr.DAO.UpdateNSRecord(c.Request.Context(), &record.Record, record); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (nr *NSRecordRouter) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	if err := nr.DAO.DeleteNSRecord(c.Request.Context(), uint(id)); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, nil)
}
