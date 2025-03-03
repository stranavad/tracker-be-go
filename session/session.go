package session

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"tracker/db"
	"tracker/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Service struct {
	types.ServiceConfig
}

func (service *Service) getSessionById(sessionId uint) (db.Session, error){
	var foundSession db.Session
	if err := service.DB.Where("id = ?", sessionId).First(&foundSession).Error; err != nil {
		return foundSession, err
	}

	return foundSession, nil
}

func (service *Service) ResetSessionTracker(c *gin.Context){
	var request ResetSessionTrackerDto

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()});
		return
	}

	if err := service.DB.Where("session_id = ?", request.SessionID).Where("tracker_id = ?", request.TrackerID).Delete(&db.Record{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
	}
}

func (service *Service) UpdateSession(c *gin.Context){
	sessionIdRaw, err := strconv.ParseUint(c.Param("sessionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid session ID param"})
		return
	}
	sessionId := uint(sessionIdRaw)

	var request StartSessionDto
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()});
		return
	}


	foundSession, err := service.getSessionById(sessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound, gin.H{"message": "Session not found"})
			return
		}

		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	foundSession.Name = request.Name
	if err := service.DB.Save(&foundSession).Error; err != nil {
		println("Failed to update session")
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, foundSession)
}

func (service *Service) StartSession(c *gin.Context){
	var request StartSessionDto
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()});
		return
	}

	var foundActiveSession db.Session
	err := service.DB.Model(&db.Session{}).Where("end_time IS NULL").First(&foundActiveSession).Error

	// If there's previous session, end it
	if err == nil {
		currentTime := time.Now()
		foundActiveSession.EndTime = &currentTime

		if err := service.DB.Save(&foundActiveSession).Error; err != nil {
			println("Failed to stop previous session")
			println(err.Error())

			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

			return
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound){
		println("Checking for ended session failed")
		println(err.Error())

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})

		return
	}


	sessionToCreate := db.Session {
		Name: request.Name,
		StartTime: time.Now(),
	}

	if err := service.DB.Save(&sessionToCreate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create new session"})
		return
	}

	if err := service.DB.Find(&sessionToCreate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create new session"})
		return
	}

	c.JSON(http.StatusOK, sessionToCreate)
}

func (service *Service) StopSession(c *gin.Context){
	sessionIdRaw, err := strconv.ParseUint(c.Param("sessionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid session ID param"})
		return
	}
	sessionId := uint(sessionIdRaw)

	session, err := service.getSessionById(sessionId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound, gin.H{"message": "Session not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
		return
	}

	if session.EndTime != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "This session is already stopped"})
		return
	}

	currentTime := time.Now()
	session.EndTime = &currentTime

	if err := service.DB.Save(&session).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to stop session"})
		return
	}

	updatedSession, err := service.getSessionById(session.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to stop session"})
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, updatedSession)
}

func (service *Service) GetSessionById(c *gin.Context){
	sessionIdRaw, err := strconv.ParseUint(c.Param("sessionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid session ID param"})
		return
	}
	sessionId := uint(sessionIdRaw)

	foundSession, err := service.getSessionById(sessionId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound, gin.H{"message": "Session not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var records []db.Record
	if err := service.DB.Where("session_id = ?", sessionId).Find(&records).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	var trackers []db.Tracker
	if err := service.DB.Order("id asc").Find(&trackers).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	trackerMap := make(map[string][]db.Record)

	for _, value := range records {
		trackerMap[value.TrackerID] = append(trackerMap[value.TrackerID], value)
	}

	trackerResponses := make([]db.TrackerResponse, len(trackers))
	for i, tracker := range trackers {
		trackerRecords := trackerMap[tracker.ID]

		var firstRecord *db.Record
		var lastRecord *db.Record

		if len(trackerRecords) > 0 {
			firstRecord = &trackerRecords[0]
			lastRecord = &trackerRecords[len(trackerRecords) - 1]
		}

		trackerRecordsResponse := make([]db.RecordResponse, len(trackerRecords))
		for index, value := range trackerRecords {
			trackerRecordsResponse[index] = value.ToResponse()
		}

		trackerResponses[i] = db.TrackerResponse{
			Tracker: tracker,
			LastRecord: lastRecord,
			FirstRecord: firstRecord,
			Records: trackerRecordsResponse,
		}
	}

	sessionResponse := db.SessionResponse{
		Session: foundSession,
		Trackers: trackerResponses,
	}

	c.JSON(http.StatusOK, sessionResponse)
}

func (service *Service) ListSessions(c *gin.Context){
	var foundSessions []db.Session

	if err := service.DB.Order("start_time desc").Find(&foundSessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
		return
	}

	c.JSON(http.StatusOK, foundSessions)
}
