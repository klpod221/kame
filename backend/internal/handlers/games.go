package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"kame/backend/internal/database"
)

type GamesHandler struct {
	db *database.DB
}

func NewGamesHandler(db *database.DB) *GamesHandler {
	return &GamesHandler{db: db}
}

func (h *GamesHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	genre := r.URL.Query().Get("genre")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	games, total, err := h.db.ListGames(q, genre, page, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve games")
		return
	}

	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"games": games,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *GamesHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing game ID")
		return
	}

	game, err := h.db.GetGameByID(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load game")
		return
	}
	if game == nil {
		respondError(w, http.StatusNotFound, "Game not found")
		return
	}

	respondJSON(w, http.StatusOK, game)
}

func (h *GamesHandler) ListGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.db.ListDistinctGenres()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load genres")
		return
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{"genres": genres})
}

type createGameRequest struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Genres        []string `json:"genres"`
	CoverURL      string   `json:"cover_url"`
	BackgroundURL string   `json:"background_url"`
	DriveFileID   string   `json:"drive_file_id"`
	FileSize      int64    `json:"file_size"`
	ReleaseYear   int      `json:"release_year"`
	Rating        float64  `json:"rating"`
}

func (h *GamesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	req.DriveFileID = strings.TrimSpace(req.DriveFileID)

	if req.Title == "" {
		respondError(w, http.StatusBadRequest, "Title is required")
		return
	}
	if req.DriveFileID == "" {
		respondError(w, http.StatusBadRequest, "Google Drive File ID is required")
		return
	}

	// Normalize Google Drive link if user pasted a full URL instead of just ID
	if strings.Contains(req.DriveFileID, "drive.google.com") {
		req.DriveFileID = extractDriveID(req.DriveFileID)
	}

	game := &database.Game{
		ID:            uuid.New().String(),
		Title:         req.Title,
		Description:   req.Description,
		Genres:        req.Genres,
		CoverURL:      req.CoverURL,
		BackgroundURL: req.BackgroundURL,
		DriveFileID:   req.DriveFileID,
		FileSize:      req.FileSize,
		ReleaseYear:   req.ReleaseYear,
		Rating:        req.Rating,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := h.db.CreateGame(game); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create game")
		return
	}

	respondJSON(w, http.StatusCreated, game)
}

func (h *GamesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing game ID")
		return
	}

	if err := h.db.DeleteGame(id); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete game")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Game deleted successfully"})
}

func extractDriveID(raw string) string {
	// e.g. https://drive.google.com/file/d/1a2b3c4d5e/view?usp=sharing
	parts := strings.Split(raw, "/d/")
	if len(parts) > 1 {
		sub := parts[1]
		if idx := strings.Index(sub, "/"); idx != -1 {
			return sub[:idx]
		}
		if idx := strings.Index(sub, "?"); idx != -1 {
			return sub[:idx]
		}
		return sub
	}
	// e.g. https://drive.google.com/open?id=1a2b3c4d5e
	if strings.Contains(raw, "id=") {
		parts = strings.Split(raw, "id=")
		if len(parts) > 1 {
			sub := parts[1]
			if idx := strings.Index(sub, "&"); idx != -1 {
				return sub[:idx]
			}
			return sub
		}
	}
	return raw
}
