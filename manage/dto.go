package manage

type CreateGroupDto struct {
	Name string `json:"name"`
}

type CreateTeamDto struct {
	Name    string `json:"name"`
	GroupID uint   `json:"groupId"`
}

type UpdateGroupDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateTeamDto struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
