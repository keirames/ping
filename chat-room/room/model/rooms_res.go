package roommodel

type JoinRoomRes struct {
	ID string `json:"id"`
}

type RoomsRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PaginateRoomsRes struct {
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Data  []RoomsRes `json:"data"`
}

type SendMessageRes struct {
	ID string `json:"id"`
}

type JoinRoomReq struct {
	RoomID string `json:"roomId" validate:"required"`
}
