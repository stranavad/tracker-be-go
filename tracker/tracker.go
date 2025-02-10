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
	identifier := c.Param("identifier")
	var record db.Record

	service.DB.Order("created_at desc").Where("identifier = ?", identifier).First(&record)

	c.JSON(http.StatusOK, record.ToResponseRecord())
}

func(service *Service) GetAllRecords(c *gin.Context){
	identifier := c.Param("identifier")
	var records []db.Record
	service.DB.Order("created_at asc").Where("identifier = ?", identifier).Find(&records)

	responseRecords := make([]db.ResponseRecord, len(records))
	for index,record := range records {
		responseRecords[index] = record.ToResponseRecord()
	}

	c.JSON(http.StatusOK, responseRecords)
}

func(service *Service) GetTrackers(c *gin.Context){
	var records []db.Record
	service.DB.Distinct("identifier").Find(&records)
	var identifiers []string
	for _, value := range records {
		identifiers = append(identifiers, value.Identifier)
	}

	c.JSON(http.StatusOK, identifiers)
}
