package database

import (
	"context"
	"fmt"
	"main/config"
	"os"

	"github.com/jackc/pgx/v5"
)

// var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
// var Conn *sqlx.DB
var Conn *pgx.Conn

func Connect() {
	conn, err :=
		pgx.
			Connect(context.Background(), config.C.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	Conn = conn
}
