package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"kame/backend/internal/config"
	"kame/backend/internal/database"
)

type AdminHandler struct {
	db  *database.DB
	cfg *config.Config
}

func NewAdminHandler(db *database.DB, cfg *config.Config) *AdminHandler {
	return &AdminHandler{
		db:  db,
		cfg: cfg,
	}
}

func (h *AdminHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.db.GetStats()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve system stats")
		return
	}

	stats["whitelist_enabled"] = len(h.cfg.WhitelistEmails) > 0
	stats["has_rawg"] = h.cfg.RAWGAPIKey != ""
	stats["has_rawg_key"] = h.cfg.RAWGAPIKey != ""
	stats["has_igdb"] = h.cfg.TwitchClientID != "" && h.cfg.TwitchClientSecret != ""
	stats["has_igdb_key"] = stats["has_igdb"]

	respondJSON(w, http.StatusOK, stats)
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.db.ListUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to load users")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"users": users})
}

type updateUserRoleRequest struct {
	Role string `json:"role"`
}

func (h *AdminHandler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		respondError(w, http.StatusBadRequest, "Missing user email")
		return
	}

	var req updateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	req.Role = strings.ToLower(strings.TrimSpace(req.Role))
	if req.Role != "admin" && req.Role != "member" {
		respondError(w, http.StatusBadRequest, "Role must be either 'admin' or 'member'")
		return
	}

	if err := h.db.UpdateUserRole(email, req.Role); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update user role")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "User role updated successfully",
		"email":   email,
		"role":    req.Role,
	})
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	if email == "" {
		respondError(w, http.StatusBadRequest, "Missing user email")
		return
	}

	// Prevent deleting the primary admin defined in env
	if len(h.cfg.AdminEmails) > 0 && h.cfg.IsAdmin(email) {
		respondError(w, http.StatusForbidden, "Cannot delete an administrator defined in system environment")
		return
	}

	if err := h.db.DeleteUser(email); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "User removed successfully",
		"email":   email,
	})
}
