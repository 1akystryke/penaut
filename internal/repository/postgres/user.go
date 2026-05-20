package postgres

import (
	"context"

	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateUser(ctx context.Context, name string, email string, passwordHash string) (model.User, error) {
	var u model.User
	err := s.db.QueryRow(ctx,
		`INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, name, email, password_hash,coalesce(picture_path, '') as picture_path`,
		name, email, passwordHash,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.PicturePath)
	return u, err
}

func (s *Store) ListUsers(ctx context.Context) ([]model.User, error) {
	rows, err := s.db.Query(ctx, `SELECT id, name, email,'penis' as password_hash,coalesce(picture_path, '') as picture_path FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
}

func (s *Store) CheckPassword(ctx context.Context, email string, passwordHash string) (*model.User, error) {
	//rows, err := s.db.Query(ctx, `SELECT id, name, email FROM users ORDER BY created_at DESC`)
	rows, err := s.db.Query(ctx,
		`Select id, name, email,password_hash,coalesce(picture_path, '') as picture_path FROM users as u
			where email = $1 and password_hash = $2 
			LIMIT 1`,
		email, passwordHash,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userModel, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
	return &userModel, err

}

func (s *Store) GetUserbyToken(ctx context.Context, token string) (*model.User, error) {
	//rows, err := s.db.Query(ctx, `SELECT id, name, email FROM users ORDER BY created_at DESC`)

	rows, err := s.db.Query(ctx,
		`select u.id, u.name, u.email,'penis' as password_hash,coalesce(u.picture_path, '') as   picture_path
		FROM users as u left join tokens as t
		on t.user_id = u.id
		where t."token" = $1
		limit 1
		`,
		token,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userModel, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
	return &userModel, err

}
func (s *Store) GetUserbyID(ctx context.Context, userID string) (*model.User, error) {
	//rows, err := s.db.Query(ctx, `SELECT id, name, email FROM users ORDER BY created_at DESC`)

	rows, err := s.db.Query(ctx,
		`select u.id, u.name, u.email,'penis' as password_hash,coalesce(picture_path, '') as  picture_path
		FROM users as u
		where u."id" = $1
		limit 1
		`,
		userID,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userModel, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.User])
	return &userModel, err

}
