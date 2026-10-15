package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"kame/backend/internal/database"
	"kame/backend/internal/scraper"
)

type RescrapeJobStatus struct {
	Status      string     `json:"status"` // "idle", "running", "completed", "cancelled", "failed"
	Total       int        `json:"total"`
	Processed   int        `json:"processed"`
	Success     int        `json:"success"`
	Failed      int        `json:"failed"`
	CurrentGame string     `json:"current_game"`
	StartedAt   time.Time  `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	LastError   string     `json:"last_error,omitempty"`
}

type RescrapeManager struct {
	mu       sync.RWMutex
	status   RescrapeJobStatus
	cancelFn context.CancelFunc
	db       *database.DB
	scraper  scraper.Provider
}

func NewRescrapeManager(db *database.DB, sc scraper.Provider) *RescrapeManager {
	return &RescrapeManager{
		db:      db,
		scraper: sc,
		status: RescrapeJobStatus{
			Status: "idle",
		},
	}
}

func (m *RescrapeManager) GetStatus(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	respondJSON(w, http.StatusOK, m.status)
}

func (m *RescrapeManager) TriggerAll(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	if m.status.Status == "running" {
		m.mu.Unlock()
		respondError(w, http.StatusConflict, "A metadata re-scrape job is already running")
		return
	}

	games, err := m.db.GetAllGameBasicInfos()
	if err != nil {
		m.mu.Unlock()
		respondError(w, http.StatusInternalServerError, "Failed to load games list from database")
		return
	}

	if len(games) == 0 {
		m.mu.Unlock()
		respondError(w, http.StatusBadRequest, "No games in vault to re-scrape")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.cancelFn = cancel

	now := time.Now()
	m.status = RescrapeJobStatus{
		Status:      "running",
		Total:       len(games),
		Processed:   0,
		Success:     0,
		Failed:      0,
		CurrentGame: "",
		StartedAt:   now,
		FinishedAt:  nil,
		LastError:   "",
	}
	m.mu.Unlock()

	// Launch background worker goroutine
	go m.runWorker(ctx, games)

	respondJSON(w, http.StatusAccepted, map[string]interface{}{
		"message": "Re-scrape job started successfully in background",
		"total":   len(games),
	})
}

func (m *RescrapeManager) Cancel(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status.Status != "running" {
		respondError(w, http.StatusBadRequest, "No active re-scrape job to cancel")
		return
	}

	if m.cancelFn != nil {
		m.cancelFn()
	}
	m.status.Status = "cancelled"
	now := time.Now()
	m.status.FinishedAt = &now

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Re-scrape job cancelled",
	})
}

func (m *RescrapeManager) runWorker(ctx context.Context, games []database.GameBasic) {
	log.Printf("[RESCRAPE] Background worker started for %d games", len(games))

	for _, g := range games {
		select {
		case <-ctx.Done():
			log.Printf("[RESCRAPE] Worker cancelled by user")
			m.mu.Lock()
			m.status.Status = "cancelled"
			now := time.Now()
			m.status.FinishedAt = &now
			m.mu.Unlock()
			return
		default:
		}

		if g.IsLocked {
			log.Printf("[RESCRAPE] Skipping locked/untracked game %q", g.Title)
			m.mu.Lock()
			m.status.Processed++
			m.mu.Unlock()
			continue
		}

		m.mu.Lock()
		m.status.CurrentGame = g.Title
		m.mu.Unlock()

		// Call scraper (Steam -> RAWG -> IGDB)
		candidates, err := m.scraper.Search(ctx, g.Title)
		if err != nil || len(candidates) == 0 {
			log.Printf("[RESCRAPE] Failed to find metadata for %q: %v", g.Title, err)
			m.mu.Lock()
			m.status.Failed++
			m.status.Processed++
			m.mu.Unlock()
		} else {
			best := candidates[0]
			// Update DB record
			updateErr := m.db.UpdateGameMetadata(
				g.ID,
				best.Title,
				best.Description,
				best.Genres,
				best.CoverURL,
				best.BackgroundURL,
				best.ReleaseYear,
				best.Rating,
			)

			m.mu.Lock()
			if updateErr != nil {
				log.Printf("[RESCRAPE] Failed to update DB for %q: %v", g.Title, updateErr)
				m.status.Failed++
			} else {
				log.Printf("[RESCRAPE] Successfully re-scraped %q from %s", best.Title, best.Source)
				m.status.Success++
			}
			m.status.Processed++
			m.mu.Unlock()
		}

		// Pacing delay (500ms) to respect rate limits and keep server responsive
		select {
		case <-ctx.Done():
			m.mu.Lock()
			m.status.Status = "cancelled"
			now := time.Now()
			m.status.FinishedAt = &now
			m.mu.Unlock()
			return
		case <-time.After(500 * time.Millisecond):
		}
	}

	m.mu.Lock()
	m.status.Status = "completed"
	m.status.CurrentGame = ""
	now := time.Now()
	m.status.FinishedAt = &now
	m.mu.Unlock()

	log.Printf("[RESCRAPE] Completed re-scraping %d games (Success: %d, Failed: %d)", len(games), m.status.Success, m.status.Failed)
}

func (m *RescrapeManager) ToggleLock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing game ID")
		return
	}

	var payload struct {
		IsLocked bool `json:"is_locked"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := m.db.ToggleGameLock(id, payload.IsLocked); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update game lock status")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":        id,
		"is_locked": payload.IsLocked,
		"message":   "Lock status updated successfully",
	})
}

func (m *RescrapeManager) RescrapeSingle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "Missing game ID")
		return
	}

	game, err := m.db.GetGameByID(id)
	if err != nil || game == nil {
		respondError(w, http.StatusNotFound, "Game not found")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 25*time.Second)
	defer cancel()

	candidates, err := m.scraper.Search(ctx, game.Title)
	if err != nil || len(candidates) == 0 {
		respondError(w, http.StatusNotFound, "Failed to scrape metadata for game")
		return
	}

	best := candidates[0]
	if err := m.db.UpdateGameMetadata(
		game.ID,
		best.Title,
		best.Description,
		best.Genres,
		best.CoverURL,
		best.BackgroundURL,
		best.ReleaseYear,
		best.Rating,
	); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save scraped metadata to database")
		return
	}

	updated, _ := m.db.GetGameByID(id)
	respondJSON(w, http.StatusOK, updated)
}
