package database

import (
	"context"
	"fmt"
	"main/config"
	"main/query"
	"os"

	"github.com/jackc/pgx/v5"
)

// var Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
// var Conn *sqlx.DB
var Conn *pgx.Conn
var Queries *query.Queries

func Connect() {
	conn, err :=
		pgx.
			Connect(context.Background(), config.C.DatabaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	Queries = query.New(conn)
	Conn = conn
	fmt.Println("DB connected!")
}
