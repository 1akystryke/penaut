package postgres

import (
	"context"

	"messenger/internal/model"

	"github.com/jackc/pgx/v5"
)

func (s *Store) GetAttachmentsbyPostID(ctx context.Context, postID string) ([]model.Attachment, error) {
	rows, err := s.db.Query(ctx, `select * from attachments where post_id = $1`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[model.Attachment])
}
