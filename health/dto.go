package health

type SaveHealthDto struct {
	DeviceID string `json:"deviceId"`
	Voltage float32 `json:"voltage"`
}


type UpdateDeviceDto struct {
	ID string `json:"id"`
	Name string `json:"name"`
}
