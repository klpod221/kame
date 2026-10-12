package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

func New(dbPath string) (*DB, error) {
	dir := filepath.Dir(dbPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create db directory %s: %w", dir, err)
		}
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Pragmas for high concurrency & robustness
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA foreign_keys = ON;",
	}
	for _, p := range pragmas {
		if _, err := conn.Exec(p); err != nil {
			log.Printf("[DB] Warning executing pragma %q: %v", p, err)
		}
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) migrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS games (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT,
			genres TEXT,
			cover_url TEXT,
			background_url TEXT,
			drive_file_id TEXT NOT NULL,
			file_size INTEGER DEFAULT 0,
			release_year INTEGER DEFAULT 0,
			rating REAL DEFAULT 0.0,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);`,
		`CREATE INDEX IF NOT EXISTS idx_games_title ON games(title);`,
		`CREATE INDEX IF NOT EXISTS idx_games_created_at ON games(created_at DESC);`,
		`CREATE TABLE IF NOT EXISTS users (
			email TEXT PRIMARY KEY,
			name TEXT,
			picture TEXT,
			role TEXT,
			last_login_at DATETIME NOT NULL
		);`,
	}

	for _, q := range queries {
		if _, err := db.conn.Exec(q); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) ListGames(search, genre string, page, limit int) ([]Game, int, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit

	var whereClauses []string
	var args []interface{}

	if strings.TrimSpace(search) != "" {
		whereClauses = append(whereClauses, "(title LIKE ? OR description LIKE ?)")
		pattern := "%" + strings.TrimSpace(search) + "%"
		args = append(args, pattern, pattern)
	}

	if strings.TrimSpace(genre) != "" {
		whereClauses = append(whereClauses, "genres LIKE ?")
		args = append(args, "%"+strings.TrimSpace(genre)+"%")
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM games %s", whereSQL)
	var total int
	if err := db.conn.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, title, description, genres, cover_url, background_url, drive_file_id, file_size, release_year, rating, created_at, updated_at
		FROM games
		%s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereSQL)

	queryArgs := append(args, limit, offset)
	rows, err := db.conn.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var g Game
		var genresRaw string
		var createdAtStr, updatedAtStr string
		err := rows.Scan(
			&g.ID,
			&g.Title,
			&g.Description,
			&genresRaw,
			&g.CoverURL,
			&g.BackgroundURL,
			&g.DriveFileID,
			&g.FileSize,
			&g.ReleaseYear,
			&g.Rating,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, 0, err
		}
		if genresRaw != "" {
			_ = json.Unmarshal([]byte(genresRaw), &g.Genres)
		}
		if g.Genres == nil {
			g.Genres = []string{}
		}
		g.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		g.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
		games = append(games, g)
	}

	if games == nil {
		games = []Game{}
	}

	return games, total, nil
}

func (db *DB) GetGameByID(id string) (*Game, error) {
	row := db.conn.QueryRow(`
		SELECT id, title, description, genres, cover_url, background_url, drive_file_id, file_size, release_year, rating, created_at, updated_at
		FROM games
		WHERE id = ?
	`, id)

	var g Game
	var genresRaw string
	var createdAtStr, updatedAtStr string
	err := row.Scan(
		&g.ID,
		&g.Title,
		&g.Description,
		&genresRaw,
		&g.CoverURL,
		&g.BackgroundURL,
		&g.DriveFileID,
		&g.FileSize,
		&g.ReleaseYear,
		&g.Rating,
		&createdAtStr,
		&updatedAtStr,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if genresRaw != "" {
		_ = json.Unmarshal([]byte(genresRaw), &g.Genres)
	}
	if g.Genres == nil {
		g.Genres = []string{}
	}
	g.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	g.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAtStr)
	return &g, nil
}

func (db *DB) CreateGame(g *Game) error {
	genresJSON, _ := json.Marshal(g.Genres)
	now := time.Now()
	if g.CreatedAt.IsZero() {
		g.CreatedAt = now
	}
	g.UpdatedAt = now

	_, err := db.conn.Exec(`
		INSERT INTO games (id, title, description, genres, cover_url, background_url, drive_file_id, file_size, release_year, rating, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		g.ID,
		g.Title,
		g.Description,
		string(genresJSON),
		g.CoverURL,
		g.BackgroundURL,
		g.DriveFileID,
		g.FileSize,
		g.ReleaseYear,
		g.Rating,
		g.CreatedAt.Format(time.RFC3339),
		g.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (db *DB) DeleteGame(id string) error {
	_, err := db.conn.Exec("DELETE FROM games WHERE id = ?", id)
	return err
}

func (db *DB) ListDistinctGenres() ([]string, error) {
	rows, err := db.conn.Query("SELECT genres FROM games WHERE genres IS NOT NULL AND genres != ''")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	genreMap := make(map[string]struct{})
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err == nil && raw != "" {
			var list []string
			if err := json.Unmarshal([]byte(raw), &list); err == nil {
				for _, g := range list {
					trimmed := strings.TrimSpace(g)
					if trimmed != "" {
						genreMap[trimmed] = struct{}{}
					}
				}
			}
		}
	}

	var results []string
	for g := range genreMap {
		results = append(results, g)
	}
	return results, nil
}

func (db *DB) UpsertUser(u *User) error {
	_, err := db.conn.Exec(`
		INSERT INTO users (email, name, picture, role, last_login_at)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT(email) DO UPDATE SET
			name = excluded.name,
			picture = excluded.picture,
			role = excluded.role,
			last_login_at = excluded.last_login_at
	`, u.Email, u.Name, u.Picture, u.Role, u.LastLoginAt.Format(time.RFC3339))
	return err
}

func (db *DB) GetUser(email string) (*User, error) {
	row := db.conn.QueryRow("SELECT email, name, picture, role, last_login_at FROM users WHERE email = ?", email)
	var u User
	var lastLoginStr string
	err := row.Scan(&u.Email, &u.Name, &u.Picture, &u.Role, &lastLoginStr)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	u.LastLoginAt, _ = time.Parse(time.RFC3339, lastLoginStr)
	return &u, nil
}
