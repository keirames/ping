package rooms

type RoomsRepository interface {
}

type roomsRepository struct{}

func (rr *roomsRepository) IsRoomExist(roomID int64) {
	// q, args, err :=
	// 	rr.Psql.
	// 		Select("id").
	// 		From("chat_rooms as cr").
	// 		Where(sq.Eq{"cr.id": roomID}).
	// 		ToSql()
	// if err != nil {
	// 	logger.L.Error().Err(err).Msg("Fail to create sql")
	// 	return false, err
	// }

	// var r int64
	// err = db.Conn.Get(&r, q, args...)
	// if err == sql.ErrNoRows {
	// 	return false, nil
	// }
	// if err != nil {
	// 	logger.L.Error().Err(err).Msg("Fail to query")
	// 	return false, err
	// }

	// return true, nil
}
