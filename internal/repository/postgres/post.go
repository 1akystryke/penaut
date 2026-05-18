package postgres

import (
	"context"

	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreatePost(ctx context.Context, channelID string, text string, userID string) (model.Post, error) {
	var p model.Post
	err := s.db.QueryRow(ctx,
		`WITH inserted AS (
			INSERT INTO posts (text, channel_id, author)
			VALUES ($1, $2, $3)
			RETURNING id, text, channel_id, author, created_at
		)
		SELECT i.*, u.name
		FROM inserted i
		left JOIN users u ON i.author = u.id`,
		text, channelID, userID,
	).Scan(&p.ID, &p.Text, &p.ChannelID, &p.UserID, &p.CreatedAt, &p.UserName)
	return p, err
}

func (s *Store) ListPosts(ctx context.Context, channelID string) ([]model.Post, error) {
	rows, err := s.db.Query(ctx,
		`SELECT p.id, p.text, p.channel_id, p.author,u.name as author_name, p.created_at
		 FROM posts as p
		 left join users as u
		 	on p.author = u.id
		 WHERE channel_id = $1
		 ORDER BY created_at`,
		channelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Post])
}
