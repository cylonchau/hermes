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
	viewIDStr := c.Query("view_id")
	var viewID *int64
	if viewIDStr != "" {
		id, _ := strconv.ParseInt(viewIDStr, 10, 64)
		viewID = &id
	}

	records, err := ar.DAO.ListARecords(c.Request.Context(), viewID)
	if err != nil {
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
