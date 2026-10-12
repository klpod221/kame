package config

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DBPath             string
	GoogleClientID     string
	WhitelistEmails    []string
	AdminEmails        []string
	JWTSecret          []byte
	RAWGAPIKey         string
	TwitchClientID     string
	TwitchClientSecret string
}

func Load() *Config {
	// Try loading .env if it exists
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "data/catalog.db")
	googleClientID := getEnv("GOOGLE_CLIENT_ID", "")
	rawgKey := getEnv("RAWG_API_KEY", "")
	twitchClientID := getEnv("TWITCH_CLIENT_ID", "")
	twitchSecret := getEnv("TWITCH_CLIENT_SECRET", "")

	whitelistRaw := getEnv("WHITELIST_EMAILS", "")
	adminRaw := getEnv("ADMIN_EMAILS", "")

	jwtSecretStr := getEnv("JWT_SECRET", "")
	var jwtSecret []byte
	if jwtSecretStr != "" {
		jwtSecret = []byte(jwtSecretStr)
	} else {
		// Generate random 32-byte secret if not provided
		randomBytes := make([]byte, 32)
		if _, err := rand.Read(randomBytes); err != nil {
			log.Fatalf("failed to generate random jwt secret: %v", err)
		}
		jwtSecret = []byte(hex.EncodeToString(randomBytes))
		log.Println("[CONFIG] Generated random ephemeral JWT_SECRET for session security")
	}

	return &Config{
		Port:               port,
		DBPath:             dbPath,
		GoogleClientID:     googleClientID,
		WhitelistEmails:    splitAndTrim(whitelistRaw),
		AdminEmails:        splitAndTrim(adminRaw),
		JWTSecret:          jwtSecret,
		RAWGAPIKey:         rawgKey,
		TwitchClientID:     twitchClientID,
		TwitchClientSecret: twitchSecret,
	}
}

func (c *Config) IsWhitelisted(email string) bool {
	if len(c.WhitelistEmails) == 0 {
		return true // If empty whitelist, allow all authenticated users
	}
	email = strings.ToLower(strings.TrimSpace(email))
	for _, w := range c.WhitelistEmails {
		if strings.ToLower(w) == email {
			return true
		}
	}
	return false
}

func (c *Config) IsAdmin(email string) bool {
	if len(c.AdminEmails) == 0 {
		// If no admin explicitly set, all whitelisted users have admin role
		return true
	}
	email = strings.ToLower(strings.TrimSpace(email))
	for _, a := range c.AdminEmails {
		if strings.ToLower(a) == email {
			return true
		}
	}
	return false
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

func splitAndTrim(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.ToLower(strings.TrimSpace(p))
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
