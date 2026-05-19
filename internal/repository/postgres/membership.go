package postgres

import (
	"context"

	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) AddMember(ctx context.Context, userID string, channelID string) (model.Membership, error) {
	var m model.Membership
	err := s.db.QueryRow(ctx,
		`INSERT INTO memberships (user_id, channel_id)
		 VALUES ($1, $2)
		 ON CONFLICT (user_id, channel_id) DO UPDATE SET user_id = EXCLUDED.user_id
		 RETURNING user_id, channel_id`,
		userID, channelID,
	).Scan(&m.UserID, &m.ChannelID)
	return m, err
}

func (s *Store) ListMembers(ctx context.Context, channelID string) ([]model.User, error) {
	rows, err := s.db.Query(ctx,
		`SELECT u.id, u.name, u.email , 'penis' as password_hash
		 FROM users u
		 JOIN memberships m ON m.user_id = u.id
		 WHERE m.channel_id = $1
		 ORDER BY u.name`,
		channelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
}
