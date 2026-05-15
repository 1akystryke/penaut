package postgres

import (
	"context"
	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) MakeToken(ctx context.Context, userID string, token string) (model.Token, error) {
	var p model.Token
	err := s.db.QueryRow(ctx,
		`INSERT INTO tokens (token, user_id)
		 VALUES ($1, $2)
		 RETURNING token,user_id, created_at`,
		token, userID,
	).Scan(&p.Token, &p.UserID, &p.CreatedAt)
	return p, err
}

func (s *Store) GetToken(ctx context.Context, token string) (*model.Token, error) {
	rows, err := s.db.Query(ctx,
		`SELECT t.*
		 FROM users as u
		 left join tokens as t
			on t.user_id=u.id 
		where t.token = $1
		 LIMIT 1`,
		token,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	TokenModel, ifnil := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.Token])
	return &TokenModel, ifnil
}
