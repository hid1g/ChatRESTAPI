package databaseCmd

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateUserDB(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS users(
	id SERIAL PRIMARY KEY,
	name VARCHAR(25) NOT NULL,
	created_at TIMESTAMP NOT NULL
	
	)
	`
	_, err := conn.Exec(ctx, sqlQuery)
	return err

}
