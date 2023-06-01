package db

import "github.com/jmoiron/sqlx"

var Conn *sqlx.DB

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
