<div align="center">
    <h1>🎮 Kame: Modern Game Vault & Direct CDN Downloader</h1>
    <p>A sleek, lightweight All-in-One Game Vault with Jellyfin-inspired dark cinema aesthetics, Google Drive high-speed direct CDN streaming, and automated metadata scraping built with Go + Vue 3.</p>
    <img src="https://img.shields.io/github/last-commit/klpod221/kame?style=for-the-badge&color=74c7ec&labelColor=111827" alt="Last Commit">
    <img src="https://img.shields.io/github/stars/klpod221/kame?style=for-the-badge&color=facc15&labelColor=111827" alt="GitHub Stars">
    <img src="https://img.shields.io/github/repo-size/klpod221/kame?style=for-the-badge&color=a78bfa&labelColor=111827" alt="Repo Size">
    <img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge&color=34d399&labelColor=111827" alt="License">
    <br/>
    <a href="https://github.com/sponsors/klpod221">
        <img src="https://img.shields.io/badge/Sponsor-GitHub-ea4aaa?style=for-the-badge&logo=github&labelColor=111827" alt="GitHub Sponsors">
    </a>
    <a href="https://www.buymeacoffee.com/klpod221">
        <img src="https://img.shields.io/badge/Buy%20Me%20a%20Coffee-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black&labelColor=111827" alt="Buy Me a Coffee">
    </a>
</div>

## 📝 Description

**Kame** is a self-contained, high-performance Game Vault and CDN streaming downloader designed for small private groups and home lab enthusiasts. Combining a sleek, Jellyfin/Steam-inspired dark cinema interface with a rock-solid single-binary Go backend, Kame allows you to browse, search, manage, and download your game collection effortlessly.

Large game archives (>10GB) hosted on Google Drive are streamed directly via client-side OAuth2 tokens with automatic abuse confirmation—bypassing Google's antivirus scan warning pages and consuming **zero server bandwidth**. Packaged entirely into a **single minimal Docker image (~25MB)** with an embedded CGO-free SQLite database.

## 🚀 Table Of Content

- [📝 Description](#-description)
- [🚀 Table Of Content](#-table-of-content)
- [✨ Features](#-features)
  - [🎮 Game Vault & Cinema UI](#-game-vault--cinema-ui)
  - [⚡ High-Speed Direct CDN Streaming](#-high-speed-direct-cdn-streaming)
  - [🤖 Automated Metadata Pipeline](#-automated-metadata-pipeline)
  - [🛡️ Security & Access Control](#️-security--access-control)
  - [📦 Single-Binary All-in-One Architecture](#-single-binary-all-in-one-architecture)
- [📦 Installation Guide](#-installation-guide)
  - [Docker Compose (Recommended)](#docker-compose-recommended)
  - [Docker CLI](#docker-cli)
- [⚙️ Configuration (.env)](#️-configuration-env)
- [🚀 Development](#-development)
  - [Prerequisites](#prerequisites)
  - [Installation & Local Run](#installation--local-run)
  - [Project Structure](#project-structure)
  - [Key Technologies](#key-technologies)
- [🔒 Security Considerations](#-security-considerations)
- [🤝 Contributing](#-contributing)
- [📝 License](#-license)
- [👤 Author](#-author)
- [🙏 Acknowledgments](#-acknowledgments)
- [📮 Support](#-support)
- [⭐ Star History](#-star-history)

---

## ✨ Features

### 🎮 Game Vault & Cinema UI
- **Obsidian Dark Aesthetic:** Jellyfin/Steam-inspired glassmorphism, dynamic color highlights, and cinematic poster cards.
- **Instant Search & Filtering:** Debounced real-time title search (`Ctrl+K`), quick genre chips, and responsive grid layouts.
- **Rich Game Details:** High-resolution hero backdrops, system requirements (minimum & recommended specs), developer & publisher details, and release dates.
- **Admin Management Panel:** Dedicated management view to quickly scrape, insert, edit, or remove games, as well as view system status.

### ⚡ High-Speed Direct CDN Streaming
- **Bypasses Google Drive Virus Warnings:** Automatically appends `acknowledgeAbuse=true&access_token={token}` to stream files >10GB directly without manual confirmation redirects.
- **Zero Server Bandwidth:** Downloads stream directly from Google Cloud CDN to the client browser at maximum speed.
- **cURL Command Generator:** Generate one-click authenticated cURL / wget download scripts for headless Linux downloading.

### 🤖 Automated Metadata Pipeline
- **Smart Title Sanitization:** Automatically cleans scene tags (`[FitGirl]`, `v1.0`, `CRACK`, `CODEX`, `.iso`, etc.) into clean search keywords.
- **Multi-Source Scraping:** Primary scraping via **RAWG API** with automatic fallback to **IGDB / Twitch API**.
- **Automated Asset Ingestion:** Fetches covers, screenshots, ratings, genres, and release years in seconds.

### 🛡️ Security & Access Control
- **Google OAuth2 Authentication:** Secure login powered by Google Identity Services (GIS).
- **Backend Whitelist Verification:** Strict email-based access control. Only authorized Google accounts can view the vault.
- **Role-Based Permissions:** Granular admin rights (`ADMIN_EMAILS`) for sensitive management endpoints.

### 📦 Single-Binary All-in-One Architecture
- **Embedded Frontend:** The Vue 3 SPA is compiled and embedded directly into the Go executable via `//go:embed`.
- **Pure Go SQLite:** Powered by `modernc.org/sqlite` (CGO-free, WAL mode, busy timeout). Zero external database container required.
- **Featherweight Footprint:** The entire production runner image is approximately **~25MB** on Alpine Linux.

---

## 📦 Installation Guide

### Docker Compose (Recommended)

1. Clone the repository and navigate to the project directory:
```bash
git clone https://github.com/klpod221/kame.git
cd kame
```

2. Copy and configure the environment variables:
```bash
cp .env.example .env
nano .env # or vim .env
```

3. Launch the container:
```bash
docker compose up -d
```

4. Access the web interface at `http://localhost:8324`.

---

### Docker CLI

You can also run Kame directly using `docker run`:

```bash
docker run -d \
  --name kame_vault \
  --restart unless-stopped \
  -p 8324:8324 \
  -v $(pwd)/data:/data \
  -e PORT=8324 \
  -e DB_PATH=/data/catalog.db \
  -e GOOGLE_CLIENT_ID="your-client-id.apps.googleusercontent.com" \
  -e WHITELIST_EMAILS="user1@gmail.com,admin@gmail.com" \
  -e ADMIN_EMAILS="admin@gmail.com" \
  -e RAWG_API_KEY="your-rawg-key" \
  kame:latest
```

---

## ⚙️ Configuration (.env)

| Environment Variable | Description | Default / Example |
| :--- | :--- | :--- |
| `PORT` | HTTP port the application listens on | `8324` |
| `DB_PATH` | Path to the SQLite database file | `/data/catalog.db` (Docker) or `data/catalog.db` (Local) |
| `GOOGLE_CLIENT_ID` | Google OAuth2 Client ID (Web Application) | `your-id.apps.googleusercontent.com` |
| `WHITELIST_EMAILS` | Comma-separated list of allowed user Gmail addresses | `user1@gmail.com,admin@gmail.com` |
| `ADMIN_EMAILS` | Comma-separated list of Admin Gmail addresses | `admin@gmail.com` |
| `JWT_SECRET` | Secret key for signing session tokens (auto-generated if empty) | `super-secret-random-32-chars` |
| `RAWG_API_KEY` | RAWG Video Games Database API Key | Free at [rawg.io/apidocs](https://rawg.io/apidocs) |
| `TWITCH_CLIENT_ID` | Twitch / IGDB Client ID (Optional fallback) | Registered at [dev.twitch.tv](https://dev.twitch.tv/console) |
| `TWITCH_CLIENT_SECRET` | Twitch / IGDB Client Secret (Optional fallback) | Registered at [dev.twitch.tv](https://dev.twitch.tv/console) |

---

## 🚀 Development

### Prerequisites
- **Go** (v1.24 or higher)
- **Node.js** (v22 or higher) & **pnpm**
- **Docker** & **Docker Compose** (optional, for containerization)

### Installation & Local Run

1. **Clone the repository:**
```bash
git clone https://github.com/klpod221/kame.git
cd kame
```

2. **Frontend Setup:**
```bash
cd frontend
pnpm install
pnpm dev
```
The frontend dev server runs at `http://localhost:5173` and proxies API requests to `http://localhost:8324`.

3. **Backend Setup:**
```bash
cd ../backend
cp ../.env.example .env
PORT=8324 DB_PATH=../data/catalog.db go run cmd/server/main.go
```

4. **Production All-in-One Build (via Makefile):**
```bash
make build
# Binary created at ./bin/kame
./bin/kame
```

---

### Project Structure

```text
kame/
├── backend/
│   ├── cmd/server/
│   │   ├── main.go               # Application entry point, Chi router & embedded SPA server
│   │   └── dist/                 # Embedded frontend production build output
│   ├── internal/
│   │   ├── config/config.go      # Environment configuration & access helpers
│   │   ├── database/             # SQLite init (WAL, pragmas, schema migrations & models)
│   │   ├── handlers/             # REST API handlers (Auth, Games, Scraper, Config)
│   │   ├── middleware/           # Auth JWT & Role-Based authorization guards
│   │   └── scraper/              # RAWG, IGDB (Twitch) & Composite Metadata Scrapers
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/           # Reusable UI components & layouts
│   │   ├── composables/          # Reactive logic (useAuth, useGames)
│   │   ├── stores/               # Pinia state stores (auth, games, admin)
│   │   ├── views/                # Views (CatalogView, GameDetailView, AdminView, LoginView)
│   │   ├── types/                # TypeScript interface definitions
│   │   └── style.css             # Tailwind CSS v4 styling
│   ├── vite.config.ts            # Vite bundler & Tailwind v4 configuration
│   └── package.json
├── Dockerfile                    # Multi-stage Dockerfile (Node -> Go -> Alpine ~25MB)
├── docker-compose.yml            # Single-service Docker Compose deployment
├── Makefile                      # Build automation shortcuts
├── LICENSE                       # MIT License
└── README.md
```

---

### Key Technologies

| Layer | Technology | Purpose |
| :--- | :--- | :--- |
| **Frontend Framework** | Vue 3 (Composition API) + TypeScript | Reactive, type-safe user interface |
| **Styling** | Tailwind CSS v4 + Lucide Icons | Ultra-fast dark cinema aesthetics |
| **State Management** | Pinia | Global reactive store management |
| **Backend Server** | Go 1.24 + Chi Router | High-concurrency, lightweight REST API |
| **SPA Embedding** | Go `//go:embed` | Bundles frontend static files directly into the Go binary |
| **Database** | SQLite (`modernc.org/sqlite`) | CGO-free, embedded transactional storage with WAL mode |
| **Scraper** | RAWG API + IGDB / Twitch API | Automatic title sanitization & game metadata fetching |
| **Authentication** | Google Identity Services (GIS) + JWT | Secure OAuth2 token validation & role whitelisting |

---

## 🔒 Security Considerations

- **Zero Drive Bandwidth on Host:** The server never downloads or re-hosts Drive files. All downloads occur via client-side scoped tokens directly from Google's CDN.
- **Email Whitelist:** Only explicitly whitelisted email addresses can access the application.
- **Admin Role Separation:** Destructive operations (adding, updating, deleting games, triggering scrapers) require dedicated admin privileges.
- **Stateless JWT Sessions:** Access tokens are signed using HMAC-SHA256 with strict expiration timeouts.
- **Embedded CGO-free Database:** Zero network ports exposed for the database, minimizing external attack vectors.

---

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the [MIT License](LICENSE) - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

**Bùi Thanh Xuân (klpod221)**

- Website: [klpod221.com](https://klpod221.com)
- GitHub: [@klpod221](https://github.com/klpod221)
- Email: [klpod221@gmail.com](mailto:klpod221@gmail.com)

---

## 🙏 Acknowledgments

- **Vue.js** & **Vite** - For the reactive and rapid frontend ecosystem.
- **Go Chi** - For the lightweight, idiomatic HTTP router.
- **RAWG** & **IGDB** - For comprehensive video game metadata APIs.
- **Google Identity Services** - For seamless OAuth2 and Google Drive integration.
- **Lucide Icons** - For the beautiful, clean icon library.

---

## 📮 Support

If you encounter any issues or have questions:

1. Check existing [Issues](https://github.com/klpod221/kame/issues)
2. Create a new issue with detailed reproduction steps
3. Contact via email: [klpod221@gmail.com](mailto:klpod221@gmail.com)

---

## ⭐ Star History

<a href="https://www.star-history.com/?repos=klpod221%2Fkame&type=date&legend=bottom-right">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=klpod221/kame&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=klpod221/kame&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=klpod221/kame&type=date&legend=top-left" />
 </picture>
</a>

---

<div align="center">
    <p>Made with ❤️ by klpod221</p>
    <p>⭐ Star this repository if you find it helpful!</p>
</div>
