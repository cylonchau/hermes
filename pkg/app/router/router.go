package router

import (
	v1 "github.com/cylonchau/hermes/pkg/app/api/v1"
	"github.com/cylonchau/hermes/pkg/dao/rdb"
	"github.com/cylonchau/hermes/pkg/model"
	"github.com/gin-gonic/gin"
)

func RegisteredRouter(e *gin.Engine) {
	// Health check
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Initialize DAOs
	zoneDAO := rdb.NewZoneDAO(model.DB)
	recordDAO := rdb.NewRecordDAO(model.DB)
	viewDAO := rdb.NewViewDAO(model.DB)

	// API V1 Group
	v1Group := e.Group("/api/v1")
	{
		// Base Resources
		zoneH := &v1.ZoneRouter{DAO: zoneDAO}
		zoneGroup := v1Group.Group("/zones")
		{
			zoneGroup.GET("", zoneH.List)
			zoneGroup.POST("", zoneH.Create)
			zoneGroup.GET("/:id", zoneH.Get)
			zoneGroup.PUT("/:id", zoneH.Update)
			zoneGroup.DELETE("/:id", zoneH.Delete)
		}

		recordBaseH := &v1.RecordRouter{DAO: recordDAO}
		recordBaseGroup := v1Group.Group("/records")
		{
			recordBaseGroup.GET("", recordBaseH.List)
			recordBaseGroup.POST("", recordBaseH.Create)
			recordBaseGroup.GET("/:id", recordBaseH.Get)
			recordBaseGroup.PUT("/:id", recordBaseH.Update)
			recordBaseGroup.DELETE("/:id", recordBaseH.Delete)
		}

		viewH := &v1.ViewRouter{DAO: viewDAO}
		viewGroup := v1Group.Group("/views")
		{
			viewGroup.GET("", viewH.List)
			viewGroup.POST("", viewH.Create)
			viewGroup.GET("/:id", viewH.Get)
			viewGroup.PUT("/:id", viewH.Update)
			viewGroup.DELETE("/:id", viewH.Delete)
		}

		// Specific Record Types
		aH := &v1.ARecordRouter{DAO: recordDAO}
		aGroup := v1Group.Group("/records/a")
		{
			aGroup.GET("", aH.List)
			aGroup.POST("", aH.Create)
			aGroup.GET("/:id", aH.Get)
			aGroup.PUT("/:id", aH.Update)
			aGroup.DELETE("/:id", aH.Delete)
		}

		aaaaH := &v1.AAAARecordRouter{DAO: recordDAO}
		aaaaGroup := v1Group.Group("/records/aaaa")
		{
			aaaaGroup.GET("", aaaaH.List)
			aaaaGroup.POST("", aaaaH.Create)
			aaaaGroup.GET("/:id", aaaaH.Get)
			aaaaGroup.PUT("/:id", aaaaH.Update)
			aaaaGroup.DELETE("/:id", aaaaH.Delete)
		}

		cnameH := &v1.CNAMERecordRouter{DAO: recordDAO}
		cnameGroup := v1Group.Group("/records/cname")
		{
			cnameGroup.GET("", cnameH.List)
			cnameGroup.POST("", cnameH.Create)
			cnameGroup.GET("/:id", cnameH.Get)
			cnameGroup.PUT("/:id", cnameH.Update)
			cnameGroup.DELETE("/:id", cnameH.Delete)
		}

		mxH := &v1.MXRecordRouter{DAO: recordDAO}
		mxGroup := v1Group.Group("/records/mx")
		{
			mxGroup.GET("", mxH.List)
			mxGroup.POST("", mxH.Create)
			mxGroup.GET("/:id", mxH.Get)
			mxGroup.PUT("/:id", mxH.Update)
			mxGroup.DELETE("/:id", mxH.Delete)
		}

		nsH := &v1.NSRecordRouter{DAO: recordDAO}
		nsGroup := v1Group.Group("/records/ns")
		{
			nsGroup.GET("", nsH.List)
			nsGroup.POST("", nsH.Create)
			nsGroup.GET("/:id", nsH.Get)
			nsGroup.PUT("/:id", nsH.Update)
			nsGroup.DELETE("/:id", nsH.Delete)
		}

		soaH := &v1.SOARecordRouter{DAO: recordDAO}
		soaGroup := v1Group.Group("/records/soa")
		{
			soaGroup.GET("", soaH.List)
			soaGroup.POST("", soaH.Create)
			soaGroup.GET("/:id", soaH.Get)
			soaGroup.PUT("/:id", soaH.Update)
			soaGroup.DELETE("/:id", soaH.Delete)
		}

		srvH := &v1.SRVRecordRouter{DAO: recordDAO}
		srvGroup := v1Group.Group("/records/srv")
		{
			srvGroup.GET("", srvH.List)
			srvGroup.POST("", srvH.Create)
			srvGroup.GET("/:id", srvH.Get)
			srvGroup.PUT("/:id", srvH.Update)
			srvGroup.DELETE("/:id", srvH.Delete)
		}

		txtH := &v1.TXTRecordRouter{DAO: recordDAO}
		txtGroup := v1Group.Group("/records/txt")
		{
			txtGroup.GET("", txtH.List)
			txtGroup.POST("", txtH.Create)
			txtGroup.GET("/:id", txtH.Get)
			txtGroup.PUT("/:id", txtH.Update)
			txtGroup.DELETE("/:id", txtH.Delete)
		}

		caaH := &v1.CAARecordRouter{DAO: recordDAO}
		caaGroup := v1Group.Group("/records/caa")
		{
			caaGroup.GET("", caaH.List)
			caaGroup.POST("", caaH.Create)
			caaGroup.GET("/:id", caaH.Get)
			caaGroup.PUT("/:id", caaH.Update)
			caaGroup.DELETE("/:id", caaH.Delete)
		}
	}
}
