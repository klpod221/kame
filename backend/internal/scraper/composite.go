package scraper

import (
	"context"
	"log"

	"kame/backend/internal/database"
)

type CompositeProvider struct {
	steam *SteamProvider
	rawg  *RAWGProvider
	igdb  *IGDBProvider
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
		steam: NewSteam(),
		rawg:  rawg,
		igdb:  igdb,
	}
}

func (c *CompositeProvider) Name() string {
	return "composite"
}

func (c *CompositeProvider) Search(ctx context.Context, query string) ([]database.GameCandidate, error) {
	var candidates []database.GameCandidate
	var err error

	// 1. Fetch rich metadata (Rating, Released, Genres, Description, Clean title) from RAWG
	if c.rawg != nil {
		candidates, err = c.rawg.Search(ctx, query)
		if err != nil {
			log.Printf("[SCRAPER] RAWG search failed for query %q: %v", query, err)
		}
	}

	// 2. Fallback to IGDB if RAWG has no results
	if len(candidates) == 0 && c.igdb != nil {
		candidates, err = c.igdb.Search(ctx, query)
		if err != nil {
			log.Printf("[SCRAPER] IGDB search failed for query %q: %v", query, err)
		}
	}

	if len(candidates) == 0 {
		return []database.GameCandidate{}, nil
	}

	// 3. Enrich cover art with official Steam 600x900 vertical capsule & hero backdrop
	if c.steam != nil {
		for i := range candidates {
			steamCover, steamHero, ok := c.steam.FetchSteamAssets(ctx, candidates[i].Title)
			if ok && steamCover != "" {
				candidates[i].CoverURL = steamCover
				if steamHero != "" {
					candidates[i].BackgroundURL = steamHero
				}
				log.Printf("[SCRAPER] Enriched %q with verified Steam 600x900 cover", candidates[i].Title)
			}
		}
	}

	return candidates, nil
}
