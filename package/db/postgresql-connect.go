package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

func Connect() (*pgx.Conn, error) {
	connStr := os.Getenv("POSTGRES_CONN")
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	log.Println("Connected to PostgreSQL")
	return conn, nil
}
