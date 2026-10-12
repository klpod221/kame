package database

import (
	"time"
)

type Game struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Genres        []string  `json:"genres"`
	CoverURL      string    `json:"cover_url"`
	BackgroundURL string    `json:"background_url"`
	DriveFileID   string    `json:"drive_file_id"`
	FileSize      int64     `json:"file_size"`
	ReleaseYear   int       `json:"release_year"`
	Rating        float64   `json:"rating"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type User struct {
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Picture     string    `json:"picture"`
	Role        string    `json:"role"` // "admin" or "user"
	LastLoginAt time.Time `json:"last_login_at"`
}

type GameCandidate struct {
	Source        string   `json:"source"` // "rawg" or "igdb"
	SourceID      string   `json:"source_id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Genres        []string `json:"genres"`
	CoverURL      string   `json:"cover_url"`
	BackgroundURL string   `json:"background_url"`
	ReleaseYear   int      `json:"release_year"`
	Rating        float64  `json:"rating"`
}
