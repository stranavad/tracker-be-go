package tracker

import (
	"tracker/types"

	"github.com/gin-gonic/gin"
)


func RegisterRoutes(r *gin.Engine, config types.ServiceConfig){
	service := Service{config}

	trackerGroup := r.Group("/tracker")
	trackerGroup.POST("/latest", service.SaveRecord)
	trackerGroup.GET("/latest/:tracker", service.GetLatestRecord)
	trackerGroup.GET("/all/:tracker", service.GetAllRecords)
	trackerGroup.GET("/trackers", service.GetTrackers)
}
