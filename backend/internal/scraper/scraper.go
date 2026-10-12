package scraper

import (
	"context"
	"regexp"
	"strings"

	"kame/backend/internal/database"
)

type Provider interface {
	Name() string
	Search(ctx context.Context, query string) ([]database.GameCandidate, error)
}

var cleanRegexes = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\[.*?\]`),                          // [FitGirl Repack], [CODEX], etc.
	regexp.MustCompile(`(?i)\(.*?\)`),                          // (v1.0), (Build 12345), etc.
	regexp.MustCompile(`(?i)\b(v\d+(\.\d+)*|repack|crack|cracked|codex|flt|skidrow|tenoke|elamigos|dodi|multi\d*)\b`),
	regexp.MustCompile(`(?i)\.(zip|rar|7z|iso|exe|tar|gz)$`),   // file extensions
	regexp.MustCompile(`[._\-]`),                               // replace dots, underscores, dashes with space
	regexp.MustCompile(`\s+`),                                  // multiple spaces
}

// CleanGameTitle removes common scene/release tags and cleanups filename to search keywords
func CleanGameTitle(raw string) string {
	cleaned := raw
	for _, r := range cleanRegexes {
		cleaned = r.ReplaceAllString(cleaned, " ")
	}
	return strings.TrimSpace(cleaned)
}
