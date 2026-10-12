package scraper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"kame/backend/internal/database"
)

type IGDBProvider struct {
	clientID     string
	clientSecret string
	client       *http.Client

	mu          sync.RWMutex
	accessToken string
	tokenExpiry time.Time
}

func NewIGDB(clientID, clientSecret string) *IGDBProvider {
	return &IGDBProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (i *IGDBProvider) Name() string {
	return "igdb"
}

func (i *IGDBProvider) getAccessToken(ctx context.Context) (string, error) {
	i.mu.RLock()
	if i.accessToken != "" && time.Now().Before(i.tokenExpiry.Add(-time.Minute)) {
		token := i.accessToken
		i.mu.RUnlock()
		return token, nil
	}
	i.mu.RUnlock()

	i.mu.Lock()
	defer i.mu.Unlock()

	// Double check after acquire write lock
	if i.accessToken != "" && time.Now().Before(i.tokenExpiry.Add(-time.Minute)) {
		return i.accessToken, nil
	}

	form := url.Values{}
	form.Set("client_id", i.clientID)
	form.Set("client_secret", i.clientSecret)
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://id.twitch.tv/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := i.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch twitch token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("twitch auth failed (%d): %s", resp.StatusCode, string(body))
	}

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	i.accessToken = res.AccessToken
	i.tokenExpiry = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)
	return i.accessToken, nil
}

type igdbGame struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Summary     string  `json:"summary"`
	Rating      float64 `json:"rating"`
	FirstRel    int64   `json:"first_release_date"`
	Cover       struct {
		ImageID string `json:"image_id"`
	} `json:"cover"`
	Screenshots []struct {
		ImageID string `json:"image_id"`
	} `json:"screenshots"`
	Genres []struct {
		Name string `json:"name"`
	} `json:"genres"`
}

func (i *IGDBProvider) Search(ctx context.Context, query string) ([]database.GameCandidate, error) {
	if i.clientID == "" || i.clientSecret == "" {
		return nil, fmt.Errorf("TWITCH_CLIENT_ID or TWITCH_CLIENT_SECRET not configured")
	}

	token, err := i.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	cleanedQuery := CleanGameTitle(query)
	if cleanedQuery == "" {
		cleanedQuery = query
	}

	// Apicalypse query syntax
	bodyStr := fmt.Sprintf(`search "%s"; fields name, summary, rating, first_release_date, cover.image_id, screenshots.image_id, genres.name; limit 5;`, cleanedQuery)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.igdb.com/v4/games", strings.NewReader(bodyStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", i.clientID)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "text/plain")

	resp, err := i.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("igdb request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("igdb error (%d): %s", resp.StatusCode, string(b))
	}

	var results []igdbGame
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode igdb response: %w", err)
	}

	var candidates []database.GameCandidate
	for _, item := range results {
		var genres []string
		for _, g := range item.Genres {
			genres = append(genres, g.Name)
		}

		releaseYear := 0
		if item.FirstRel > 0 {
			t := time.Unix(item.FirstRel, 0)
			releaseYear = t.Year()
		}

		coverURL := ""
		if item.Cover.ImageID != "" {
			coverURL = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg", item.Cover.ImageID)
		}

		bgURL := coverURL
		if len(item.Screenshots) > 0 && item.Screenshots[0].ImageID != "" {
			bgURL = fmt.Sprintf("https://images.igdb.com/igdb/image/upload/t_1080p/%s.jpg", item.Screenshots[0].ImageID)
		}

		candidates = append(candidates, database.GameCandidate{
			Source:        "igdb",
			SourceID:      strconv.Itoa(item.ID),
			Title:         item.Name,
			Description:   item.Summary,
			Genres:        genres,
			CoverURL:      coverURL,
			BackgroundURL: bgURL,
			ReleaseYear:   releaseYear,
			Rating:        item.Rating / 20.0, // Scale 100 to 5-star
		})
	}

	return candidates, nil
}
