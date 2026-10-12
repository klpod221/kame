package handlers

import (
	"net/http"

	"kame/backend/internal/config"
)

type ConfigHandler struct {
	cfg *config.Config
}

func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{cfg: cfg}
}

func (h *ConfigHandler) Get(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"google_client_id":   h.cfg.GoogleClientID,
		"whitelist_enabled":  len(h.cfg.WhitelistEmails) > 0,
		"has_rawg_key":       h.cfg.RAWGAPIKey != "",
		"has_igdb":           h.cfg.TwitchClientID != "" && h.cfg.TwitchClientSecret != "",
	})
}
