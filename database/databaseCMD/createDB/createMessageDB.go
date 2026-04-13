package databaseCmd

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateMessageDB(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS messages(
	id SERIAL PRIMARY KEY,
	sended_from INTEGER NOT NULL,
	sended_to INTEGER NOT NULL, 
	text_message VARCHAR NOT NULL, 
	is_sended BOOLEAN NOT NULL,
	sended_time TIMESTAMP, 
	is_read BOOLEAN NOT NULL
	)
	`
	_, err := conn.Exec(ctx, sqlQuery)
	return err
}
