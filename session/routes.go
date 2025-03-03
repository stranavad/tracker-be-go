package session

import (
	"tracker/types"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, config types.ServiceConfig){
	service := Service{config}

	sessionGroup := r.Group("/session")
	sessionGroup.POST("", service.StartSession)
	sessionGroup.POST("/stop/:sessionId", service.StopSession)
	sessionGroup.POST("/reset-tracker", service.ResetSessionTracker)
	sessionGroup.PUT("/:sessionId", service.UpdateSession)
	sessionGroup.GET("/:sessionId", service.GetSessionById)
	sessionGroup.GET("/list", service.ListSessions)
}
