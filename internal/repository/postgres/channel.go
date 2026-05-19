package postgres

import (
	"context"

	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateChannel(ctx context.Context, channelType string, channelName string) (model.Channel, error) {
	var ch model.Channel
	err := s.db.QueryRow(ctx,
		`INSERT INTO channels (type,name) VALUES ($1,$2) RETURNING id, type,name`,
		channelType, channelName,
	).Scan(&ch.ID, &ch.Type, &ch.Name)
	return ch, err
}

func (s *Store) ListChannels(ctx context.Context, userId string) ([]model.Channel, error) {
	rows, err := s.db.Query(ctx, `SELECT c.id, max(c.type) as type,max(c.name) as name 
				FROM channels as c
				left join memberships m
					on m.channel_id = c.id
				where m.user_id  = $1
				or
				c."type" = 'public'
				group by c.id`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Channel])
}
