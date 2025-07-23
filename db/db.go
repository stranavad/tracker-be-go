package db

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormModelDefault struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Session struct {
	GormModelDefault
	Name      string     `json:"name"`
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
}

type Group struct {
	GormModelDefault
	Name  string `json:"name"`
	Teams []Team `json:"teams"`
}

type Team struct {
	GormModelDefault
	Name    string `json:"name"`
	GroupID uint   `json:"groupId"`
}

type Tracker struct {
	ID    string `gorm:"primarykey" json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Device struct {
	ID    string `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type DeviceResponse struct {
	ID      string                 `json:"id"`
	Name    string                 `json:"name"`
	Color   string                 `json:"color"`
	Records []DeviceHealthResponse `json:"records"`
}

func (tracker *Tracker) ToResponse(healthRecords []Record) DeviceResponse {
	records := make([]DeviceHealthResponse, len(healthRecords))
	for i, record := range healthRecords {
		records[i] = record.ToResponseHealth()
	}

	return DeviceResponse{
		ID:      tracker.ID,
		Color:   tracker.Color,
		Name:    tracker.Name,
		Records: records,
	}
}

func (device *Device) ToResponse(healthRecords []DeviceHealth) DeviceResponse {
	records := make([]DeviceHealthResponse, len(healthRecords))
	for i, record := range healthRecords {
		records[i] = record.ToResponseHealth()
	}

	return DeviceResponse{
		ID:      device.ID,
		Name:    device.Name,
		Color:   device.Color,
		Records: records,
	}
}

type DeviceHealth struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	DeviceID  string    `json:"deviceId"`
	Voltage   float64   `json:"voltage"`
	Trace     string    `json:"trace"`
}

func (deviceHealth *DeviceHealth) ToResponseHealth() DeviceHealthResponse {
	return DeviceHealthResponse{
		Trace:     deviceHealth.Trace,
		Timestamp: deviceHealth.CreatedAt.Unix(),
		Voltage:   deviceHealth.Voltage,
	}
}

type Record struct {
	GormModelDefault
	Lat             float32 `json:"lat"`
	Long            float32 `json:"long"`
	TrackerID       string  `json:"trackerId"`
	SessionID       *uint   `json:"sessionId"`
	Trace           string  `json:"trace"`
	Voltage         float64 `json:"voltage"`
	DeviceTimestamp int64   `json:"timestamp"`
}

func (record *Record) ToResponseHealth() DeviceHealthResponse {
	return DeviceHealthResponse{
		Trace:     record.Trace,
		Timestamp: record.DeviceTimestamp,
		Voltage:   record.Voltage,
	}
}

type DeviceHealthResponse struct {
	Trace     string  `json:"trace"`
	Timestamp int64   `json:"timestamp"`
	Voltage   float64 `json:"voltage"`
}

type SessionResponse struct {
	Session
	Trackers []TrackerResponse `json:"trackers"`
}

type TrackerResponse struct {
	Tracker
	Records []Record `json:"records"`
}

var db *gorm.DB

func init() {
	envErr := godotenv.Load()
	if envErr != nil {
		fmt.Println("Error loading .env file")
	}
	connStr := os.Getenv("DATABASE_URL")

	fmt.Println("Connecting to DB")
	var err error
	db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		panic("Failed to connect database")
	}

	err = db.AutoMigrate(&Session{}, &Tracker{}, &Record{}, &Device{}, &DeviceHealth{}, &Group{}, &Team{})
	if err != nil {
		panic("Failed to migrate database")
	}
}

func GetDb() *gorm.DB {
	return db
}
