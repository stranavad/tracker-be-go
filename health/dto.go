package health

type SaveHealthDto struct {
	DeviceID string `json:"deviceId"`
	Voltage float64 `json:"voltage"`
	Trace string `json:"trace"`
}


type UpdateDeviceDto struct {
	ID string `json:"id"`
	Name string `json:"name"`
}
