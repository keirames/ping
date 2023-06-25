package messagemodel

import "database/sql"

type MessageEntity struct {
	ID        int64         `db:"id"`
	Content   string        `db:"content"`
	Type      string        `db:"type"`
	IsDelete  bool          `db:"is_delete"`
	ParentID  sql.NullInt64 `db:"parent_id"`
	CreatedAt string        `db:"created_at"`
	UserID    int64         `db:"user_id"`
	RoomID    int64         `db:"room_id"`
}

type MessageRes struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	IsDelete  bool   `json:"isDelete"`
	ParentID  *int64 `json:"parentId"`
	CreatedAt string `json:"createdAt"`
	UserID    int64  `json:"userId"`
	RoomID    int64  `json:"roomId"`
}

func MapMessageEntityToModel(me MessageEntity) MessageRes {
	var parentID *int64
	if me.ParentID.Valid {
		parentID = &me.ParentID.Int64
	}

	return MessageRes{
		ID:        me.ID,
		Content:   me.Content,
		Type:      me.Type,
		IsDelete:  me.IsDelete,
		ParentID:  parentID,
		CreatedAt: me.CreatedAt,
		UserID:    me.UserID,
		RoomID:    me.RoomID,
	}
}

func MapMessagesEntityToModel(messages []MessageEntity) []MessageRes {
	res := make([]MessageRes, len(messages))

	for idx, m := range messages {
		res[idx] = MapMessageEntityToModel(m)
	}

	return res
}
