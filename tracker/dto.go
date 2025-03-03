package tracker

import "tracker/db"



type SaveRecordDto struct {
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
	Rssi int16 `json:"rssi"`
	Snr int8 `json:"snr"`
	Identifier string `json:"identifier"`
}


type UpdateTrackerDto struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

func (dto *SaveRecordDto) ToModel(sessionId *uint) db.Record {
	return db.Record{
		Lat: dto.Lat,
		Long: dto.Long,
		Rssi: dto.Rssi,
		Snr: dto.Snr,
		TrackerID: dto.Identifier,
		SessionID: sessionId,
	}
}
