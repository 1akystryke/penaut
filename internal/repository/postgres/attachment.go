package postgres

import (
	"context"
	"log"
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

func (s *Store) UpdatePostAttachment(ctx context.Context, picName string, postID string, fileType string) (*model.Attachment, error) {
	log.Println("вызов записи в субд")
	var attachment model.Attachment
	err := s.db.QueryRow(ctx,
		`INSERT INTO attachments (file_path,post_id,file_type) VALUES ($1,$2,$3) RETURNING file_path, post_id,file_type`,
		picName, postID, fileType,
	).Scan(&attachment.FilePath, &attachment.PostID, &attachment.FileType)
	if err != nil {
		log.Println("ошибка записи в субд ", err)
	}
	return &attachment, err
}
