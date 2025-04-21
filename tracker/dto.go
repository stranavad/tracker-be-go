package tracker

import "tracker/db"

type GetLastRecordsDto struct {
	TrackerID string `json:"trackerId"`
	LastRecordID *uint `json:"lastRecordId"`
}

type SaveRecordDto struct {
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
	Identifier string `json:"identifier"`
	Trace string `json:"trace"`
	Timestamp int64 `json:"timestamp"`
}


type UpdateTrackerDto struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

func (dto *SaveRecordDto) ToModel(sessionId *uint) db.Record {
	return db.Record{
		Lat: dto.Lat,
		Long: dto.Long,
		TrackerID: dto.Identifier,
		SessionID: sessionId,
		Trace: dto.Trace,
		DeviceTimestamp: dto.Timestamp,
	}
}
