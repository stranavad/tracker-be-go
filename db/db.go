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
	ID        uint `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

type Session struct {
	GormModelDefault
	Name string `json:"name"`
	StartTime time.Time `json:"startTime"`
	EndTime *time.Time `json:"endTime"`
}

type Tracker struct {
	ID string `gorm:"primarykey" json:"id"`
	Name string `json:"name"`
}


type Record struct {
	GormModelDefault
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
	Rssi int16 `json:"rssi"`
	Snr int8 `json:"snr"`
	TrackerID string `json:"trackerId"`
	SessionID *uint `json:"sessionId"`
}

/* Response types */
type SessionResponse struct {
	Session
	Trackers []TrackerResponse `json:"trackers"`
}

type TrackerResponse struct {
	Tracker
	LastRecord *Record `json:"latestRecord"`
	FirstRecord *Record `json:"firstRecord"`
	Records []RecordResponse `json:"records"`
}

func(r *Record) ToResponse() RecordResponse {
	return RecordResponse{
		Lat: r.Lat,
		Long: r.Long,
	}
}

type RecordResponse struct {
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
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

	err = db.AutoMigrate(&Session{}, &Tracker{}, &Record{})
	if err != nil {
		panic("Failed to migrate database")
	}
}

func GetDb() *gorm.DB {
	return db
}
