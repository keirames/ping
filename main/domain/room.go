package domain

type Room struct {
	ID        int64  `db:"id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	DeletedAt string `db:"deleted_at"`
}

type RoomRepository interface {
	FindByID(id int64) (*Room, error)
	FindAll() (*[]Room, error)
	CreateRoom(name string, members []int64) (*Room, error)
}

type RoomUseCase interface {
	GetAll() (*[]Room, error)
	CreateRoom(input CreateRoomInput) (*Room, error)
	// DeleteRoom() (*Room, error)
	// SendMessage() (*Message, error)
	// AddMember() (*User, error)
}

type CreateRoomInput struct {
	Name      string  `json:"name" binding:"required"`
	FriendIDs []int64 `json:"friendIds" binding:"required"`
}
