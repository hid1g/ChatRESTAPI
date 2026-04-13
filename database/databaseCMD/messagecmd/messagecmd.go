package messagecmd

import (
	"chat/chat"
	"context"

	"github.com/jackc/pgx/v5"
)

func SendMessage(ctx context.Context, conn *pgx.Conn, message chat.Message) error {
	sqlQuery := `
	INSERT INTO messages (sended_from, sended_to, text_message, is_sended, sended_time, is_read)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := conn.Exec(
		ctx,
		sqlQuery,
		message.SendedFrom,
		message.SendedTo,
		message.Text,
		message.IsSended,
		message.SendedTime,
		message.IsRead,
	)
	return err
}

func GetMessageByUser(ctx context.Context, conn *pgx.Conn, id int) ([]chat.Message, error) {
	sqlQuery := `
	SELECT id, sended_from, sended_to, text_message, is_sended, sended_time, is_read
	FROM messages
	WHERE sended_from = $1
	ORDER BY id ASC
	`
	rows, err := conn.Query(ctx, sqlQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userMessage := make([]chat.Message, 0)
	for rows.Next() {
		var usMes chat.Message
		if err := rows.Scan(
			&usMes.Id,
			&usMes.SendedFrom,
			&usMes.SendedTo,
			&usMes.Text,
			&usMes.IsSended,
			&usMes.SendedTime,
			&usMes.IsRead,
		); err != nil {
			return nil, err
		}
		userMessage = append(userMessage, usMes)
	}
	return userMessage, nil
}

func DeleteMessage(ctx context.Context, conn *pgx.Conn, id int) error {
	sqlQuery := `
	DELETE FROM messages
	WHERE id = $1
	`
	_, err := conn.Exec(ctx, sqlQuery, id)
	return err
}

func MessageIsRead(ctx context.Context, conn *pgx.Conn, id int) error {
	sqlQuery := `
	UPDATE messages
	SET is_read = true
	WHERE id = $1
	`

	_, err := conn.Exec(ctx, sqlQuery, id)
	return err
}

func MessageUpdate(ctx context.Context, conn *pgx.Conn, id int, newmes string) error {
	sqlQuery := `
	UPDATE messages
	SET text_message = $1
	WHERE id = $2
	`

	_, err := conn.Exec(ctx, sqlQuery, newmes, id)
	return err
}

func GetMessagesBetweenUsers(
	ctx context.Context,
	conn *pgx.Conn,
	userId1 int,
	userId2 int,
) ([]chat.Message, error) {
	sqlQuery := `
	SELECT id, sended_from, sended_to, text_message, is_sended, sended_time, is_read
	FROM messages
	WHERE (sended_from = $1 AND sended_to = $2)
	OR (sended_from = $2 AND sended_to = $1)
	ORDER BY id ASC	
	`
	rows, err := conn.Query(ctx, sqlQuery, userId1, userId2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messagesBetweenUsers := make([]chat.Message, 0)
	for rows.Next() {
		var usMes chat.Message
		if err := rows.Scan(
			&usMes.Id,
			&usMes.SendedFrom,
			&usMes.SendedTo,
			&usMes.Text,
			&usMes.IsSended,
			&usMes.SendedTime,
			&usMes.IsRead,
		); err != nil {
			return nil, err
		}
		messagesBetweenUsers = append(messagesBetweenUsers, usMes)
	}
	return messagesBetweenUsers, nil
}
