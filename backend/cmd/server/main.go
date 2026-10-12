package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"kame/backend/internal/config"
	"kame/backend/internal/database"
	"kame/backend/internal/handlers"
	"kame/backend/internal/middleware"
	"kame/backend/internal/scraper"
)

//go:embed all:dist
var embeddedFS embed.FS

func main() {
	cfg := config.Load()

	log.Printf("[SERVER] Starting Kame Game Catalog on port %s...", cfg.Port)
	log.Printf("[CONFIG] SQLite DB Path: %s", cfg.DBPath)
	if len(cfg.WhitelistEmails) > 0 {
		log.Printf("[CONFIG] Whitelist enabled (%d emails)", len(cfg.WhitelistEmails))
	} else {
		log.Println("[CONFIG] Whitelist is OPEN (all authenticated Google accounts allowed)")
	}

	db, err := database.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("[DB] Fatal error initializing database: %v", err)
	}
	defer db.Close()

	// Scraper engine (Composite: RAWG -> IGDB fallback)
	scraperProvider := scraper.NewComposite(cfg.RAWGAPIKey, cfg.TwitchClientID, cfg.TwitchClientSecret)

	// Handlers
	authHandler := handlers.NewAuthHandler(cfg, db)
	gamesHandler := handlers.NewGamesHandler(db)
	scraperHandler := handlers.NewScraperHandler(scraperProvider)
	configHandler := handlers.NewConfigHandler(cfg)

	r := chi.NewRouter()

	// Base Middlewares
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:*", "http://127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Auth extraction middleware
	r.Use(middleware.AuthMiddleware(cfg))

	// API Routes
	r.Route("/api", func(api chi.Router) {
		api.Get("/config", configHandler.Get)

		api.Route("/auth", func(auth chi.Router) {
			auth.Post("/google", authHandler.GoogleLogin)
			auth.Get("/me", authHandler.Me)
		})

		api.Route("/games", func(games chi.Router) {
			games.Get("/", gamesHandler.List)
			games.Get("/{id}", gamesHandler.Get)
		})

		api.Get("/genres", gamesHandler.ListGenres)

		// Protected Admin Routes
		api.Route("/admin", func(admin chi.Router) {
			admin.Use(middleware.RequireAdmin)
			admin.Post("/scrape", scraperHandler.Scrape)
			admin.Post("/games", gamesHandler.Create)
			admin.Delete("/games/{id}", gamesHandler.Delete)
		})
	})

	// Static Web Assets / SPA Fallback
	distSub, err := fs.Sub(embeddedFS, "dist")
	if err != nil {
		log.Fatalf("failed to create sub filesystem: %v", err)
	}
	r.NotFound(spaHandler(distSub))

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[SERVER] Listen error: %v", err)
		}
	}()

	log.Printf("[SERVER] Ready at http://localhost:%s", cfg.Port)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("[SERVER] Shutting down gracefully...")
	_ = server.Close()
	log.Println("[SERVER] Exited.")
}

func spaHandler(staticFS fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(staticFS))

	return func(w http.ResponseWriter, r *http.Request) {
		// Do not intercept /api calls with 404
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		// Check if file exists in static FS
		f, err := staticFS.Open(path)
		if err == nil {
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routing
		indexFile, err := staticFS.Open("index.html")
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer indexFile.Close()

		stat, err := indexFile.Stat()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
	}
}
