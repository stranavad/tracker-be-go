package tracker

import "tracker/db"


type SaveRecordDto struct {
	Lat float32 `json:"lat"`
	Long float32 `json:"long"`
	Rssi int16 `json:"rssi"`
	Snr int8 `json:"snr"`
	Identifier string `json:"identifier"`
}


func (dto *SaveRecordDto) ToModel() db.Record {
	return db.Record{
		Lat: dto.Lat,
		Long: dto.Long,
		Rssi: dto.Rssi,
		Snr: dto.Snr,
		Identifier: dto.Identifier,
	}
}