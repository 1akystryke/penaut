package postgres

import (
	"context"

	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateChannel(ctx context.Context, channelType string) (model.Channel, error) {
	var ch model.Channel
	err := s.db.QueryRow(ctx,
		`INSERT INTO channels (type) VALUES ($1) RETURNING id, type`,
		channelType,
	).Scan(&ch.ID, &ch.Type)
	return ch, err
}

func (s *Store) ListChannels(ctx context.Context) ([]model.Channel, error) {
	rows, err := s.db.Query(ctx, `SELECT id, type FROM channels ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Channel])
}
