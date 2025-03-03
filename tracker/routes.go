package tracker

import (
	"tracker/types"

	"github.com/gin-gonic/gin"
)


func RegisterRoutes(r *gin.Engine, config types.ServiceConfig){
	service := Service{config}

	trackerGroup := r.Group("/tracker")
	trackerGroup.POST("/tracker", service.SaveRecord)
	trackerGroup.PUT("", service.UpdateTracker)
	trackerGroup.GET("/latest/:trackerId", service.GetLatestRecord)
	trackerGroup.GET("/all/:trackerId", service.GetAllRecords)
	trackerGroup.GET("/trackers", service.GetTrackers)
}
