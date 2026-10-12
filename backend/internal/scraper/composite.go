package scraper

import (
	"context"
	"log"

	"kame/backend/internal/database"
)

type CompositeProvider struct {
	rawg *RAWGProvider
	igdb *IGDBProvider
}

func NewComposite(rawgKey, twitchClientID, twitchSecret string) *CompositeProvider {
	var rawg *RAWGProvider
	if rawgKey != "" {
		rawg = NewRAWG(rawgKey)
	}

	var igdb *IGDBProvider
	if twitchClientID != "" && twitchSecret != "" {
		igdb = NewIGDB(twitchClientID, twitchSecret)
	}

	return &CompositeProvider{
		rawg: rawg,
		igdb: igdb,
	}
}

func (c *CompositeProvider) Name() string {
	return "composite"
}

func (c *CompositeProvider) Search(ctx context.Context, query string) ([]database.GameCandidate, error) {
	// Try RAWG first
	if c.rawg != nil {
		results, err := c.rawg.Search(ctx, query)
		if err == nil && len(results) > 0 {
			return results, nil
		}
		if err != nil {
			log.Printf("[SCRAPER] RAWG search failed for query %q: %v. Attempting fallback...", query, err)
		}
	}

	// Fallback to IGDB
	if c.igdb != nil {
		results, err := c.igdb.Search(ctx, query)
		if err == nil && len(results) > 0 {
			return results, nil
		}
		if err != nil {
			log.Printf("[SCRAPER] IGDB search failed for query %q: %v", query, err)
		}
	}

	return []database.GameCandidate{}, nil
}
