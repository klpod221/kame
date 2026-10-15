package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"kame/backend/internal/database"
)

type SteamProvider struct {
	client *http.Client
}

func NewSteam() *SteamProvider {
	return &SteamProvider{
		client: &http.Client{Timeout: 8 * time.Second},
	}
}

func (s *SteamProvider) Name() string {
	return "steam"
}

type steamSearchItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type steamSearchResponse struct {
	Total int               `json:"total"`
	Items []steamSearchItem `json:"items"`
}

var dlcBlacklistWords = []string{
	"upgrade", "deluxe edition upgrade", "dlc", "soundtrack", "expansion",
	"season pass", "content pack", "bonus pack", "artbook", "demo", "costume",
	"multiplayer pack", "vr patch",
}

// FetchSteamAssets searches Steam for a base game and verifies if library_600x900.jpg exists.
// Returns coverURL, backgroundURL, and ok (true if verified 200 OK).
func (s *SteamProvider) FetchSteamAssets(ctx context.Context, title string) (string, string, bool) {
	cleaned := CleanGameTitle(title)
	if cleaned == "" {
		cleaned = title
	}

	searchURL := fmt.Sprintf(
		"https://store.steampowered.com/api/storesearch/?term=%s&l=english&cc=US",
		url.QueryEscape(cleaned),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, searchURL, nil)
	if err != nil {
		return "", "", false
	}
	req.Header.Set("User-Agent", "KameVault/1.0 (+https://github.com/klpod221/kame)")

	resp, err := s.client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if resp != nil {
			resp.Body.Close()
		}
		return "", "", false
	}

	var searchRes steamSearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchRes)
	resp.Body.Close()
	if err != nil || len(searchRes.Items) == 0 {
		return "", "", false
	}

	// Filter out DLCs, upgrades, soundtracks to pick the true base game
	var selectedID int
	for _, item := range searchRes.Items {
		lowerName := strings.ToLower(item.Name)
		isDLC := false
		for _, bl := range dlcBlacklistWords {
			if strings.Contains(lowerName, bl) {
				isDLC = true
				break
			}
		}
		if !isDLC {
			selectedID = item.ID
			break
		}
	}

	// If all items matched blacklist, fallback to the very first item
	if selectedID == 0 {
		selectedID = searchRes.Items[0].ID
	}

	// Candidate URLs
	steamCover := fmt.Sprintf("https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/%d/library_600x900.jpg", selectedID)
	steamHero := fmt.Sprintf("https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/%d/library_hero.jpg", selectedID)

	// Verify cover image exists via HEAD request (timeout 2.5s)
	checkCtx, checkCancel := context.WithTimeout(ctx, 3*time.Second)
	defer checkCancel()

	headReq, err := http.NewRequestWithContext(checkCtx, http.MethodHead, steamCover, nil)
	if err != nil {
		return "", "", false
	}
	headReq.Header.Set("User-Agent", "KameVault/1.0 (+https://github.com/klpod221/kame)")

	headResp, err := s.client.Do(headReq)
	if err != nil {
		return "", "", false
	}
	defer headResp.Body.Close()

	if headResp.StatusCode == http.StatusOK {
		return steamCover, steamHero, true
	}

	return "", "", false
}

// Search implements the Provider interface if needed directly
func (s *SteamProvider) Search(ctx context.Context, query string) ([]database.GameCandidate, error) {
	// Not used directly in hybrid mode, returns empty
	return []database.GameCandidate{}, nil
}
