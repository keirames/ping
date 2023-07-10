package database

import (
	"main/config"
	"main/logger"

	"github.com/jmoiron/sqlx"
)

var Conn *sqlx.DB

func Connect() {
	conn, err :=
		sqlx.
			Connect(config.C.DBDriverName, config.C.DBSource)
	if err != nil {
		logger.L.Err(err).Msg("Cannot connect to database")
		panic("Cannot connect to database")
	}

	Conn = conn
}
