package model

import "time"

type User struct {
	ID           string `json:"id" db:"id"`
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
}

type Channel struct {
	ID   string `json:"id" db:"id"`
	Type string `json:"type" db:"type"`
}

type Post struct {
	ID        string    `json:"id" db:"id"`
	Text      string    `json:"text" db:"text"`
	ChannelID string    `json:"channel_id" db:"channel_id"`
	UserID    string    `json:"author" db:"author"`
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
