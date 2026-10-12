# 🎮 KAME - Game Catalog & Direct CDN Downloader

A lightweight, self-contained All-in-One Game Catalog Web Application designed for small private groups (2-3 users). Features a sleek Steam/Jellyfin-inspired dark aesthetic, high-speed Google Drive download streaming (>10GB) that bypasses virus scan warnings, packaged entirely into a **single minimal Docker image (~25MB)**.

---

## 🌟 Key Features

* **Ultra-Lightweight All-in-One Architecture:**
  * Backend written in **Go**, serving REST APIs and embedding the entire Frontend SPA into a single static binary via `//go:embed`.
  * Pure Go **SQLite engine** (`modernc.org/sqlite`, CGO-free, preconfigured with WAL mode and busy timeout). Zero external database container required.
* **Modern & Responsive UI:**
  * **Vue 3** + **TypeScript** + **Vite** + **Tailwind CSS v4** (`@tailwindcss/vite`) + **Lucide Icons** (`lucide-vue-next`).
  * 3:4 aspect-ratio game posters, high-resolution hero banners, dynamic genre filtering, and debounced instant search.
* **High-Speed Direct Download (Bypasses Google Drive Virus Scan):**
  * Integrated with **Google Identity Services (GIS)** using scope `https://www.googleapis.com/auth/drive.readonly`.
  * Client generates direct stream URLs with `acknowledgeAbuse=true&access_token={user_access_token}`:
    ```text
    https://www.googleapis.com/drive/v3/files/{drive_file_id}?alt=media&acknowledgeAbuse=true&access_token={user_access_token}
    ```
  * Files stream directly from Google's CDN to the user's browser at full bandwidth. The application server consumes **0 bytes of download bandwidth**.
* **Automated Metadata Scraping Pipeline:**
  * Automatic filename sanitization removes scene/release tags (`[FitGirl]`, `v1.0`, `CRACK`, `CODEX`, `.iso`, etc.).
  * Fetches cover art, backgrounds, genres, release year, and ratings via **RAWG API** (primary) with automatic fallback to **IGDB / Twitch API**.
* **Private Access Control:**
  * Google OAuth2 authentication paired with backend **Email Whitelist** verification. Only whitelisted Google accounts can sign in. Includes a dedicated Admin role.

---

## 📂 Project Structure

```text
kame/
├── backend/
│   ├── cmd/server/
│   │   ├── main.go               # Entry point, Chi router, embedded SPA file server fallback
│   │   └── dist/                 # Embedded Frontend static build output
│   ├── internal/
│   │   ├── config/config.go      # Environment variable loader & whitelist/admin helpers
│   │   ├── database/             # SQLite init (WAL, pragmas, migrations, models)
│   │   ├── handlers/             # Auth, Games, Scraper, Config REST handlers
│   │   ├── middleware/           # JWT verification, Whitelist & Admin guards
│   │   └── scraper/              # RAWG, IGDB (Twitch OAuth), & Composite Scraper
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/           # Navbar, GameCard, GameDetailModal, AddGameModal
│   │   ├── composables/          # useAuth (GIS SDK), useGames (Download URL generator)
│   │   ├── types/                # TypeScript interfaces
│   │   ├── App.vue
│   │   └── style.css             # Tailwind v4 import & custom styles
│   ├── vite.config.ts            # Tailwind v4 plugin + API proxy
│   └── package.json
├── Dockerfile                    # Multi-stage build (Node -> Go Static Binary -> Alpine Runner)
├── docker-compose.yml            # Single-service Docker Compose configuration + /data volume
├── .env.example                  # Environment variables template
└── Makefile                      # Automation tasks for dev & build
```

---

## ⚙️ Configuration (.env)

Copy `.env.example` to `.env`:

```bash
cp .env.example .env
```

| Environment Variable | Description | Default / Example |
| :--- | :--- | :--- |
| `PORT` | HTTP port the application listens on | `8080` |
| `DB_PATH` | Path to the SQLite database file | `/data/catalog.db` (Docker) or `data/catalog.db` (Local) |
| `GOOGLE_CLIENT_ID` | OAuth 2.0 Client ID from Google Cloud Console | `xxxx.apps.googleusercontent.com` |
| `WHITELIST_EMAILS` | Comma-separated list of whitelisted Gmail addresses | `alice@gmail.com,bob@gmail.com` |
| `ADMIN_EMAILS` | Comma-separated list of Gmail addresses with Admin role | `alice@gmail.com` |
| `JWT_SECRET` | Secret key for signing JWT sessions | Auto-generated random 32-byte key if empty |
| `RAWG_API_KEY` | Free API key from [rawg.io/apidocs](https://rawg.io/apidocs) | `abcdef...` |
| `TWITCH_CLIENT_ID` | Client ID for IGDB fallback (optional) | `...` |
| `TWITCH_CLIENT_SECRET`| Client Secret for IGDB fallback (optional) | `...` |

> **Google Cloud Console Setup:**
> 1. Visit [Google Cloud Console Credentials](https://console.cloud.google.com/apis/credentials).
> 2. Create an **OAuth 2.0 Client ID** with Application type **Web application**.
> 3. Under **Authorized JavaScript origins**, add your origin (e.g. `http://localhost:8080` or `https://kame.yourdomain.com`).
> 4. In Google Drive, manually grant **Viewer** permissions to the 2-3 whitelisted users on the folder containing the games.

---

## 🚀 Docker Deployment (Single Container)

### 1. Run with Docker Compose

```bash
# Build image and start in background
docker compose up -d --build

# Inspect logs
docker compose logs -f
```

### 2. Run Directly with Docker CLI

```bash
# Build single Docker image
docker build -t kame:latest .

# Run container with SQLite volume mount
docker run -d \
  --name kame_catalog \
  --restart unless-stopped \
  -p 8080:8080 \
  --env-file .env \
  -v $(pwd)/data:/data \
  kame:latest
```

The application is now accessible at **`http://localhost:8080`**.

---

## 🛠️ Local Development

For local development on Arch Linux / macOS / Linux:

```bash
# 1. Start Go Backend (Terminal 1)
make dev-backend

# 2. Start Vite Frontend with Hot-Reload (Terminal 2)
make dev-frontend
```

The Vite dev server runs at `http://localhost:5173` and automatically proxies all `/api/*` requests to the Go backend on port `8080`.

To build a standalone production binary locally:

```bash
make build
# Output executable: ./bin/kame (~13MB)
./bin/kame
```
