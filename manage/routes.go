package manage

import (
	"github.com/gin-gonic/gin"
	"tracker/types"
)

func RegisterRoutes(r *gin.Engine, config types.ServiceConfig) {
	service := Service{config}

	manageGroup := r.Group("/manage")
	manageGroup.GET("/all", service.GetAll)
	manageGroup.POST("/team", service.CreateTeam)
	manageGroup.POST("/group", service.CreateGroup)
	manageGroup.PUT("/team", service.UpdateTeam)
	manageGroup.PUT("/group", service.UpdateGroup)
	manageGroup.DELETE("/group/:groupId", service.DeleteGroup)
	manageGroup.DELETE("/team/:teamId", service.DeleteTeam)
}
