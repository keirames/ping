package model

type JoinRoomRes struct {
	ID string `json:"id"`
}

type RoomsRes struct {
	ID   string
	Name string
}

type PaginateRoomsRes struct {
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Data  []RoomsRes `json:"data"`
}
