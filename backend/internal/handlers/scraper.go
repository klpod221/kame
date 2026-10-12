package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"kame/backend/internal/scraper"
)

type ScraperHandler struct {
	provider scraper.Provider
}

func NewScraperHandler(provider scraper.Provider) *ScraperHandler {
	return &ScraperHandler{provider: provider}
}

type scrapeRequest struct {
	Query string `json:"query"`
}

func (h *ScraperHandler) Scrape(w http.ResponseWriter, r *http.Request) {
	var req scrapeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	cleanedQuery := scraper.CleanGameTitle(req.Query)
	if strings.TrimSpace(cleanedQuery) == "" {
		cleanedQuery = strings.TrimSpace(req.Query)
	}

	if cleanedQuery == "" {
		respondError(w, http.StatusBadRequest, "Query cannot be empty")
		return
	}

	candidates, err := h.provider.Search(r.Context(), cleanedQuery)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Metadata search failed: "+err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"query":         req.Query,
		"cleaned_query": cleanedQuery,
		"candidates":    candidates,
	})
}
