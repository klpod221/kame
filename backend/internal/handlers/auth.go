package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"kame/backend/internal/config"
	"kame/backend/internal/database"
	"kame/backend/internal/middleware"
)

type AuthHandler struct {
	cfg *config.Config
	db  *database.DB
}

func NewAuthHandler(cfg *config.Config, db *database.DB) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db}
}

type googleAuthRequest struct {
	Credential  string `json:"credential"`   // Google ID Token from GIS
	AccessToken string `json:"access_token"` // Optional Access Token
}

type googleTokenInfo struct {
	Email         string `json:"email"`
	EmailVerified any    `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Sub           string `json:"sub"`
	Aud           string `json:"aud"`
}

func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var req googleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var email, name, picture string

	if req.Credential != "" {
		// Verify Google ID Token via Google's tokeninfo endpoint
		info, err := verifyGoogleIDToken(req.Credential)
		if err != nil {
			respondError(w, http.StatusUnauthorized, fmt.Sprintf("Google token verification failed: %v", err))
			return
		}
		email = info.Email
		name = info.Name
		picture = info.Picture
	} else if req.AccessToken != "" {
		// Verify via userinfo endpoint using access_token
		info, err := verifyGoogleAccessToken(req.AccessToken)
		if err != nil {
			respondError(w, http.StatusUnauthorized, fmt.Sprintf("Google access token verification failed: %v", err))
			return
		}
		email = info.Email
		name = info.Name
		picture = info.Picture
	} else {
		respondError(w, http.StatusBadRequest, "Missing credential or access_token")
		return
	}

	if email == "" {
		respondError(w, http.StatusBadRequest, "Email not found in Google account")
		return
	}

	// Whitelist Check
	if !h.cfg.IsWhitelisted(email) {
		respondError(w, http.StatusForbidden, fmt.Sprintf("Email %s is not in the authorized whitelist", email))
		return
	}

	// Role Determination
	role := "user"
	if h.cfg.IsAdmin(email) {
		role = "admin"
	}

	user := &database.User{
		Email:       email,
		Name:        name,
		Picture:     picture,
		Role:        role,
		LastLoginAt: time.Now(),
	}

	if err := h.db.UpsertUser(user); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save user profile")
		return
	}

	// Generate JWT Session Token (Valid for 30 days)
	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	claims := middleware.Claims{
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.cfg.JWTSecret)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to sign authentication token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"token": tokenString,
		"user":  user,
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	email, _ := r.Context().Value(middleware.UserEmailKey).(string)
	if email == "" {
		respondJSON(w, http.StatusOK, map[string]interface{}{"user": nil})
		return
	}

	user, err := h.db.GetUser(email)
	if err != nil || user == nil {
		respondJSON(w, http.StatusOK, map[string]interface{}{"user": nil})
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{"user": user})
}

func verifyGoogleIDToken(idToken string) (*googleTokenInfo, error) {
	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tokeninfo returned status %d", resp.StatusCode)
	}

	var info googleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

func verifyGoogleAccessToken(accessToken string) (*googleTokenInfo, error) {
	endpoint := "https://www.googleapis.com/oauth2/v3/userinfo"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userinfo returned status %d", resp.StatusCode)
	}

	var info googleTokenInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
