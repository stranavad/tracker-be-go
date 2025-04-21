package health

import (
	"tracker/types"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, config types.ServiceConfig){
	service := Service{config}

	healthGroup := r.Group("/health")
	healthGroup.GET("/all", service.GetHealthData)
	healthGroup.POST("/save", service.SaveHealth)
	healthGroup.PUT("/device", service.UpdateDevice)
}
