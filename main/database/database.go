package database

import (
	"main/tools"

	"github.com/jmoiron/sqlx"
)

var Conn *sqlx.DB

func Connect() error {
	conn, err := sqlx.Connect(tools.Config.DBDriver, tools.Config.DBSource)
	if err != nil {
		return err
	}

	Conn = conn

	return nil
}
