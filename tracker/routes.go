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
	trackerGroup.GET("/trackers", service.GetTrackers)
	trackerGroup.GET("/health", service.GetTrackersHealth)
}
