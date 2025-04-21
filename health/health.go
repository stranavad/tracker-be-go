package health

import (
	"errors"
	"math"
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

func(service *Service) getDeviceById(deviceId string)(db.Device, error){
	var foundDevice db.Device
	if err := service.DB.Where("id = ?", deviceId).First(&foundDevice).Error; err != nil {
		return foundDevice, err
	}

	return foundDevice, nil
}

func (service *Service) UpdateDevice(c *gin.Context){
	var request UpdateDeviceDto
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	foundDevice, err := service.getDeviceById(request.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			c.JSON(http.StatusNotFound, gin.H{"message": "Device not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		println(err.Error())
		return
	}

	foundDevice.Name = request.Name
	foundDevice.Color = request.Color

	if err := service.DB.Save(&foundDevice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()});
		println(err.Error())
		return
	}
}

func(service *Service) GetHealthData(c *gin.Context){
	var devices []db.Device
	if err := service.DB.Find(&devices).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}


	// Tak tady se mi trochu nepovedlo to udelat v jedne query takze N+1 here you are
	responseDevices := make([]db.DeviceResponse, len(devices))
	for i, device := range devices {
		var deviceHealth []db.DeviceHealth

		if err := service.DB.Where("device_id = ?", device.ID).Order("created_at DESC").Limit(100).Find(&deviceHealth).Error; err != nil {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		responseDevices[i] = device.ToResponse(deviceHealth)
	}

	c.JSON(http.StatusOK, responseDevices)
}


func (service *Service) SaveHealth(c *gin.Context){
	var request SaveHealthDto

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()});
		return
	}

	if _, err := service.getDeviceById(request.DeviceID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound){
			newDevice := db.Device{
				ID: request.DeviceID,
				Name: request.DeviceID,
			}

			if err := service.DB.Create(&newDevice).Error; err != nil {
				println("Failed to create new device")
				println(err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
				return
			}
		} else {
			println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	// Round voltage
	ratio := math.Pow(10, float64(2))
	roundedVoltage :=  math.Round(request.Voltage*ratio) / ratio

	healthToCreate := db.DeviceHealth {
		DeviceID: request.DeviceID,
		CreatedAt: time.Now(),
		Voltage: roundedVoltage,
		Trace: request.Trace,
	}

	if err := service.DB.Create(&healthToCreate).Error; err != nil {
		println(err.Error())
		c.JSON(http.StatusInternalServerError, "Failed saving record")
	}
}
