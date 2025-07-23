package session

type StartSessionDto struct {
	Name string `json:"name"`
}

type ResetSessionTrackerDto struct {
	TrackerID string `json:"trackerId"`
	SessionID uint   `json:"sessionId"`
}

type UpdateTeamToTrackerDto struct {
	SessionID uint   `json:"sessionId"`
	TrackerID string `json:"trackerId"`
	TeamID    uint   `json:"teamId"`
}
