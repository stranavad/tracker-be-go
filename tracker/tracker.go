package tracker

import (
	"fmt"
	"net/http"
	"tracker/db"
	"tracker/types"

	"github.com/gin-gonic/gin"
)


type Service struct {
	types.ServiceConfig
}


func (service *Service) SaveRecord(c *gin.Context){
	var request SaveRecordDto
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
		return
	}

	recordToCreate := request.ToModel()

	service.DB.Save(&recordToCreate)
}

func (service *Service) GetLatestRecord(c *gin.Context){
	var record db.Record

	service.DB.Order("created_at desc").First(&record)

	c.JSON(http.StatusOK, record)
}

func(service *Service) GetAllRecords(c *gin.Context){
	var records []db.Record
	service.DB.Order("created_at asc").Find(&records)

	c.JSON(http.StatusOK, records)
}
