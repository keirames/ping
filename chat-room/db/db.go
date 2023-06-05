package db

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var Conn *sqlx.DB

var Psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func Connect() error {
	conn, err :=
		sqlx.
			Connect("postgres", "postgresql://postgres:password@localhost:5432/chat-room?sslmode=disable")
	if err != nil {
		return err
	}

	Conn = conn

	return nil
}
