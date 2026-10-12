package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"kame/backend/internal/database"
)

type RAWGProvider struct {
	apiKey string
	client *http.Client
}

func NewRAWG(apiKey string) *RAWGProvider {
	return &RAWGProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *RAWGProvider) Name() string {
	return "rawg"
}

type rawgSearchResponse struct {
	Results []struct {
		ID              int     `json:"id"`
		Name            string  `json:"name"`
		Released        string  `json:"released"`
		BackgroundImage string  `json:"background_image"`
		Rating          float64 `json:"rating"`
		Genres          []struct {
			Name string `json:"name"`
		} `json:"genres"`
		ShortScreenshots []struct {
			Image string `json:"image"`
		} `json:"short_screenshots"`
	} `json:"results"`
}

type rawgGameDetail struct {
	DescriptionRaw string `json:"description_raw"`
	Description    string `json:"description"`
}

func (r *RAWGProvider) Search(ctx context.Context, query string) ([]database.GameCandidate, error) {
	if r.apiKey == "" {
		return nil, fmt.Errorf("RAWG_API_KEY is not configured")
	}

	cleanedQuery := CleanGameTitle(query)
	if cleanedQuery == "" {
		cleanedQuery = query
	}

	endpoint := fmt.Sprintf("https://api.rawg.io/api/games?key=%s&search=%s&page_size=5",
		url.QueryEscape(r.apiKey),
		url.QueryEscape(cleanedQuery),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rawg request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rawg returned status %d", resp.StatusCode)
	}

	var data rawgSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode rawg response: %w", err)
	}

	var candidates []database.GameCandidate
	for _, item := range data.Results {
		var genres []string
		for _, g := range item.Genres {
			genres = append(genres, g.Name)
		}

		releaseYear := 0
		if len(item.Released) >= 4 {
			if y, err := strconv.Atoi(item.Released[:4]); err == nil {
				releaseYear = y
			}
		}

		bg := item.BackgroundImage
		cover := item.BackgroundImage
		if len(item.ShortScreenshots) > 1 && item.ShortScreenshots[1].Image != "" {
			bg = item.ShortScreenshots[1].Image
		}

		candidate := database.GameCandidate{
			Source:        "rawg",
			SourceID:      strconv.Itoa(item.ID),
			Title:         item.Name,
			Genres:        genres,
			CoverURL:      cover,
			BackgroundURL: bg,
			ReleaseYear:   releaseYear,
			Rating:        item.Rating,
		}

		// Try fetching game detail for description
		if desc, err := r.getGameDescription(ctx, item.ID); err == nil && desc != "" {
			candidate.Description = desc
		}

		candidates = append(candidates, candidate)
	}

	return candidates, nil
}

func (r *RAWGProvider) getGameDescription(ctx context.Context, id int) (string, error) {
	endpoint := fmt.Sprintf("https://api.rawg.io/api/games/%d?key=%s", id, url.QueryEscape(r.apiKey))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	var detail rawgGameDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return "", err
	}

	desc := strings.TrimSpace(detail.DescriptionRaw)
	if desc == "" {
		desc = strings.TrimSpace(detail.Description)
	}
	return desc, nil
}
