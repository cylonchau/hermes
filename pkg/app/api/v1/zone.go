package v1

import (
	"strconv"

	"github.com/cylonchau/hermes/pkg/app/api/query"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

type ZoneRouter struct {
	DAO *rdb.ZoneDAO
}

func (zr *ZoneRouter) List(c *gin.Context) {
	zones, err := zr.DAO.GetAll(c.Request.Context(), 0, 0)
	if err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, zones)
}

func (zr *ZoneRouter) Create(c *gin.Context) {
	var zone model.Zone
	if err := c.ShouldBindJSON(&zone); err != nil {
		query.BadRequest(c, err)
		return
	}
	if err := zr.DAO.Create(c.Request.Context(), &zone); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, zone)
}

func (zr *ZoneRouter) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	zone, err := zr.DAO.GetByID(c.Request.Context(), id)
	if err != nil {
		query.NotFound(c, query.ErrZoneNotFound)
		return
	}
	query.SuccessResponse(c, nil, zone)
}

func (zr *ZoneRouter) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	zone, err := zr.DAO.GetByID(c.Request.Context(), id)
	if err != nil {
		query.NotFound(c, query.ErrZoneNotFound)
		return
	}

	if err := c.ShouldBindJSON(&zone); err != nil {
		query.BadRequest(c, err)
		return
	}

	if err := zr.DAO.Update(c.Request.Context(), zone); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, zone)
}

func (zr *ZoneRouter) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := zr.DAO.Delete(c.Request.Context(), id); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, nil)
}
