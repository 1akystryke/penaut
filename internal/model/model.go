package model

import "time"

type User struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	PicturePath  string `json:"picture_path" db:"picture_path"`
}

type Channel struct {
	ID          string `json:"id" db:"id"`
	Type        string `json:"type" db:"type"`
	Name        string `json:"name" db:"name"`
	PicturePath string `json:"picture_path" db:"picture_path"`
}

type Post struct {
	ID        string    `json:"id" db:"id"`
	Text      string    `json:"text" db:"text"`
	ChannelID string    `json:"channel_id" db:"channel_id"`
	UserID    string    `json:"author" db:"author"`
	UserName  string    `json:"user_name" db:"author_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Membership struct {
	UserID    string `json:"user_id" db:"user_id"`
	ChannelID string `json:"channel_id" db:"channel_id"`
}
type Token struct {
	Token     string    `json:"token" db:"token"`
	UserID    string    `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Attachment struct {
	FilePath string `json:"file_path" db:"file_path"`
	PostID   string `json:"post_id" db:"post_id"`
	FileType string `json:"file_type" db:"file_type"`
}
