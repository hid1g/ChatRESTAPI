package databaseCMD

import (
	"chat/chat"
	"context"

	"github.com/jackc/pgx/v5"
)

func InsertUser(ctx context.Context, conn *pgx.Conn, user chat.User) error {
	sqlQuery := `
	INSERT INTO users (name, created_at)
	VALUES ($1, $2)
	`
	_, err := conn.Exec(ctx, sqlQuery, user.Name, user.CreatedAt)
	return err
}

func ListUsers(ctx context.Context, conn *pgx.Conn) ([]chat.User, error) {
	sqlQuery := `
		SELECT id, name, created_at
		FROM users
		ORDER BY id ASC
	`
	rows, err := conn.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	usersList := make([]chat.User, 0)
	for rows.Next() {
		var user chat.User
		if err := rows.Scan(&user.Id, &user.Name, &user.CreatedAt); err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}
	return usersList, nil
}

func ListUsersByName(ctx context.Context, conn *pgx.Conn, name string) ([]chat.User, error) {
	sqlQuery := `
	SELECT id, name, created_at
	FROM users
	WHERE name = $1
	`
	rows, err := conn.Query(ctx, sqlQuery, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userList := make([]chat.User, 0)
	for rows.Next() {
		var user chat.User
		if err := rows.Scan(&user.Id, &user.Name, &user.CreatedAt); err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}
	return userList, nil
}

func DeleteUser(ctx context.Context, conn *pgx.Conn, id int) error {
	sqlQuery := `
	DELETE FROM users
	WHERE id = $1
	`

	_, err := conn.Exec(ctx, sqlQuery, id)
	return err
}

func UpdateUser(ctx context.Context, conn *pgx.Conn, user chat.User) error {
	sqlQuery := `
	UPDATE users 
	SET name=$1
	WHERE id = $2
	`
	_, err := conn.Exec(ctx, sqlQuery, user.Name, user.Id)
	return err
}
