package session

type StartSessionDto struct {
	Name string `json:"name"`
}


type ResetSessionTrackerDto struct {
	TrackerID string `json:"trackerId"`
	SessionID uint `json:"sessionId"`
}
