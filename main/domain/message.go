package domain

type Message struct {
	ID        int64  `db:"id"`
	Text      string `db:"text"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	DeletedAt string `db:"deleted_at"`
	RoomID    int64  `db:"room_id"`
	UserID    int64  `db:"user_id"`
}
