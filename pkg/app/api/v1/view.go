package v1

import (
	"strconv"

	"github.com/cylonchau/hermes/pkg/app/api/query"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

type ViewRouter struct {
	DAO *rdb.ViewDAO
}

func (vr *ViewRouter) List(c *gin.Context) {
	views, err := vr.DAO.GetAll(c.Request.Context())
	if err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, views)
}

func (vr *ViewRouter) Create(c *gin.Context) {
	var view model.View
	if err := c.ShouldBindJSON(&view); err != nil {
		query.BadRequest(c, err)
		return
	}
	if err := vr.DAO.Create(c.Request.Context(), &view); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, view)
}

func (vr *ViewRouter) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	view, err := vr.DAO.GetByID(c.Request.Context(), id)
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}
	query.SuccessResponse(c, nil, view)
}

func (vr *ViewRouter) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	view, err := vr.DAO.GetByID(c.Request.Context(), id)
	if err != nil {
		query.NotFound(c, query.ErrParam)
		return
	}

	if err := c.ShouldBindJSON(&view); err != nil {
		query.BadRequest(c, err)
		return
	}

	if err := vr.DAO.Update(c.Request.Context(), view); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, view)
}

func (vr *ViewRouter) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := vr.DAO.Delete(c.Request.Context(), id); err != nil {
		query.InternalError(c, err)
		return
	}
	query.SuccessResponse(c, nil, nil)
}
