package manage

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"tracker/db"
	"tracker/types"
)

type Service struct {
	types.ServiceConfig
}

func (service *Service) GetAll(c *gin.Context) {
	var groups []db.Group
	if err := service.DB.Preload("Teams").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to load groups"})
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, groups)
}

func (service *Service) CreateGroup(c *gin.Context) {
	var request CreateGroupDto

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	groupToCreate := db.Group{
		Name: request.Name,
	}

	if err := service.DB.Save(&groupToCreate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create group"})
		println(err.Error())
		return
	}

	if err := service.DB.Find(&groupToCreate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create group"})
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, groupToCreate)
}

func (service *Service) CreateTeam(c *gin.Context) {
	var request CreateTeamDto
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var foundGroup db.Group
	if err := service.DB.First(&foundGroup, request.GroupID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Group not found"})
		return
	}

	teamToCreate := db.Team{
		Name:    request.Name,
		GroupID: foundGroup.ID,
	}

	if err := service.DB.Save(&teamToCreate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create team"})
		println(err.Error())
		return
	}

	if err := service.DB.Find(&teamToCreate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create team"})
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, teamToCreate)
}

func (service *Service) UpdateGroup(c *gin.Context) {
	var request UpdateGroupDto
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	group, err := service.getGroupById(request.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Group not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
		return
	}

	group.Name = request.Name

	if err := service.DB.Save(&group).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

func (service *Service) UpdateTeam(c *gin.Context) {
	var request UpdateTeamDto
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	team, err := service.getTeamById(request.ID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Team not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
		return
	}

	team.Name = request.Name

	if err := service.DB.Save(&team).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, team)
}

func (service *Service) DeleteGroup(c *gin.Context) {
	groupIdRaw, err := strconv.ParseUint(c.Param("groupId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid group ID param"})
		return
	}
	groupId := uint(groupIdRaw)

	group, err := service.getGroupById(groupId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Group not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
		return
	}

	if err := service.DB.Delete(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete group"})
		println(err.Error())
		return
	}
}

func (service *Service) DeleteTeam(c *gin.Context) {
	teamIdRaw, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid group ID param"})
		return
	}
	teamId := uint(teamIdRaw)

	team, err := service.getTeamById(teamId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Team not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
		return
	}

	if err := service.DB.Delete(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete team"})
		println(err.Error())
		return
	}
}

func (service *Service) getGroupById(groupId uint) (db.Group, error) {
	var foundGroup db.Group
	if err := service.DB.Where("id = ?", groupId).First(&foundGroup).Error; err != nil {
		return foundGroup, err
	}

	return foundGroup, nil
}

func (service *Service) getTeamById(teamId uint) (db.Team, error) {
	var foundTeam db.Team
	if err := service.DB.Where("id = ?", teamId).First(&foundTeam).Error; err != nil {
		return foundTeam, err
	}

	return foundTeam, nil
}
