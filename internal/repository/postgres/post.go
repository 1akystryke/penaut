package postgres

import (
	"context"

	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreatePost(ctx context.Context, channelID string, text string) (model.Post, error) {
	var p model.Post
	err := s.db.QueryRow(ctx,
		`INSERT INTO posts (text, channel_id)
		 VALUES ($1, $2)
		 RETURNING id, text, channel_id, created_at`,
		text, channelID,
	).Scan(&p.ID, &p.Text, &p.ChannelID, &p.CreatedAt)
	return p, err
}

func (s *Store) ListPosts(ctx context.Context, channelID string) ([]model.Post, error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, text, channel_id, created_at
		 FROM posts
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
