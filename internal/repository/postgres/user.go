package postgres

import (
	"context"

	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateUser(ctx context.Context, name string, email string) (model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email`,
		name, email,
	).Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}

func (s *Store) ListUsers(ctx context.Context) ([]model.User, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name, email FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
}
