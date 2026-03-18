package v1

import (
	"strconv"

	"github.com/cylonchau/hermes/pkg/app/api/query"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

type RecordRouter struct {
	DAO *rdb.RecordDAO
}

func (rr *RecordRouter) List(c *gin.Context) {
	var zoneID int64
	var viewID *int64
	if zoneIDStr := c.Query("zone_id"); zoneIDStr != "" {
		zoneID, _ = strconv.ParseInt(zoneIDStr, 10, 64)
	}
	if viewIDStr := c.Query("view_id"); viewIDStr != "" {
		id, _ := strconv.ParseInt(viewIDStr, 10, 64)
		viewID = &id
	}

	records, err := rr.DAO.ListRecords(c.Request.Context(), zoneID, viewID)
	if err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, records)
}

func (rr *RecordRouter) Create(c *gin.Context) {
	var record model.Record
	if err := c.ShouldBindJSON(&record); err != nil {
		query.BadRequest(c, err)
		return
	}
	if err := rr.DAO.CreateRecord(c.Request.Context(), &record); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (rr *RecordRouter) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	record, err := rr.DAO.GetRecordByID(c.Request.Context(), id)
	if err != nil {
		query.NotFound(c, query.ErrRecordNotFound)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (rr *RecordRouter) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	record, err := rr.DAO.GetRecordByID(c.Request.Context(), id)
	if err != nil {
		query.NotFound(c, query.ErrRecordNotFound)
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		query.BadRequest(c, err)
		return
	}

	if err := rr.DAO.UpdateRecord(c.Request.Context(), record); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, record)
}

func (rr *RecordRouter) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := rr.DAO.DeleteRecord(c.Request.Context(), id); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, nil)
}
