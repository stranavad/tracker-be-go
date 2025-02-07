package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type Record struct {
	gorm.Model
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
	Rssi int16 `json:"rssi"`
	Snr int8 `json:"snr"`
	Identifier string `json:"identifier"`
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

	err = db.AutoMigrate(&Record{})
	if err != nil {
		panic("Failed to migrate database")
	}
}

func GetDb() *gorm.DB {
	return db
}
