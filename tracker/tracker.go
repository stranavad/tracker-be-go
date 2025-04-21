package tracker

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"tracker/db"
	"tracker/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)


type Service struct {
	types.ServiceConfig
}

func (service *Service) getTrackerById(trackerId string) (db.Tracker, error) {
	var foundTracker db.Tracker
	if err := service.DB.Where("id = ?", trackerId).First(&foundTracker).Error; err != nil {
		return foundTracker, err
	}

	return foundTracker, nil
}


func (service *Service) getCurrentSession()(*db.Session){
	var foundSession db.Session
	if err := service.DB.Where("end_time IS null").First(&foundSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			// Create base session
			sessionToCreate := db.Session{
				Name: time.Now().String(),
				StartTime: time.Now(),
			}

			if err := service.DB.Create(&sessionToCreate).Error; err != nil {
				println("Failed to create dummy session")
				return nil
			}

			if err := service.DB.First(&sessionToCreate).Error; err != nil {
				println("Failed to create dummy session")
				return nil
			}

			return &sessionToCreate
		} else {
			println("Failed to get session")
			return nil
		}
	}

	return &foundSession
}

func (service *Service) UpdateTracker(c *gin.Context){
	var request UpdateTrackerDto
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	foundTracker, err := service.getTrackerById(request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound, gin.H{"message": "Tracker not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
		return
	}

	foundTracker.Name = request.Name
	foundTracker.Color = request.Color

	if err := service.DB.Save(&foundTracker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()});
		println(err.Error())
		return
	}
}

func (service *Service) SaveRecord(c *gin.Context){
	var request SaveRecordDto
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err.Error())
		return
	}

	if request.Lat == 0 && request.Long == 0 {
		println("We don't accept zero coordinates, for now")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cannot store 0,0 coords"})
		return
	}

	session := service.getCurrentSession()
	var sessionId *uint
	if session != nil {
		sessionId = &session.ID
	}

	recordToCreate := request.ToModel(sessionId)

	// Check if the tracker is already inserted
	if _, err := service.getTrackerById(recordToCreate.TrackerID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			// Create tracker
			if err := service.DB.Create(&db.Tracker{
				ID: recordToCreate.TrackerID,
				Name: recordToCreate.TrackerID,
			}).Error; err != nil {
				println(err.Error())
				c.JSON(http.StatusInternalServerError, "Failed creating tracker")
				return
			}
		} else {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, "Failed tracker checking")
			return
		}
	}


	if err := service.DB.Save(&recordToCreate).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, "Failed saving record")
	}
}


func(service *Service) GetTrackers(c *gin.Context){
	var trackers []db.Tracker

	if err := service.DB.Find(&trackers).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}

	c.JSON(http.StatusOK, trackers)
}

func(service *Service) GetTrackersHealth(c *gin.Context){
	var trackers []db.Tracker

	if err := service.DB.Find(&trackers).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return;
	}

	responseTrackers := make([]db.DeviceResponse, len(trackers))

	for i, tracker := range trackers {
		var records []db.Record

		if err := service.DB.Where("tracker_id = ?", tracker.ID).Order("device_timestamp DESC").Limit(1).Find(&records).Error; err != nil {
			println(err.Error());
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responseTrackers[i] = tracker.ToResponse(records)
	}

	c.JSON(http.StatusOK, responseTrackers)
}
