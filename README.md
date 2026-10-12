# 🎮 KAME - Game Catalog & Direct CDN Downloader

Hệ thống Web Game Catalog All-in-One tối giản dành cho nhóm nội bộ (2-3 người dùng) với giao diện dark-mode phong cách Steam/Jellyfin, tích hợp cơ chế tải file Google Drive dung lượng lớn (>10GB) không bị chặn quét virus và được đóng gói toàn bộ vào **1 Docker image duy nhất (~25MB)**.

---

## 🌟 Tính Năng Nổi Bật

* **Kiến Trúc All-in-One Siêu Nhẹ:** 
  * Backend viết bằng **Go**, phục vụ REST API và nhúng toàn bộ Frontend SPA vào 1 binary tĩnh duy nhất qua `//go:embed`.
  * Cơ sở dữ liệu **SQLite thuần Go** (`modernc.org/sqlite`, CGO-free, bật sẵn WAL mode, busy timeout). Không cần cài đặt container PostgreSQL rời.
* **Giao Diện Hiện Đại & Mượt Mà:**
  * **Vue 3** + **Vite** + **Tailwind CSS v4** (`@tailwindcss/vite`) + **Lucide Icons** (`lucide-vue-next`).
  * Poster tỷ lệ 3:4, Hero banner chất lượng cao, bộ lọc Genre động và tìm kiếm tức thì.
* **Tải File Tốc Độ Tối Đa (Bypass Virus Scan Google Drive):**
  * Tích hợp **Google Identity Services (GIS)** với scope `https://www.googleapis.com/auth/drive.readonly`.
  * Sinh URL tải trực tiếp kèm `acknowledgeAbuse=true&access_token={user_access_token}` trên Client:
    ```text
    https://www.googleapis.com/drive/v3/files/{drive_file_id}?alt=media&acknowledgeAbuse=true&access_token={user_access_token}
    ```
  * File được stream thẳng từ máy chủ CDN của Google về máy người dùng. Máy chủ ứng dụng **không tốn 1 byte băng thông nào**.
* **Tự Động Cào Metadata (Scraper Pipeline):**
  * Tự động làm sạch tên file/game (lọc các tag như `[FitGirl]`, `v1.0`, `CRACK`, `CODEX`, `.iso`,...).
  * Tìm kiếm và lấy poster, background art, thể loại, năm phát hành, điểm đánh giá qua **RAWG API** (ưu tiên) và tự động fallback sang **IGDB / Twitch API**.
* **Bảo Mật Nội Bộ:**
  * Xác thực qua Google OAuth2 kết hợp kiểm tra **Whitelist Email** ở phía Backend. Chỉ các email được khai báo trước mới có quyền truy cập. Phân quyền Admin riêng biệt.

---

## 📂 Cấu Trúc Mã Nguồn

```text
kame/
├── backend/
│   ├── cmd/server/
│   │   ├── main.go               # Entry point, Chi router, SPA file server fallback
│   │   └── dist/                 # Thư mục nhúng static assets của Frontend
│   ├── internal/
│   │   ├── config/config.go      # Đọc biến môi trường & helper whitelist/admin
│   │   ├── database/             # SQLite init (WAL, pragma, migrations, models)
│   │   ├── handlers/             # Auth, Games, Scraper, Config REST handlers
│   │   ├── middleware/           # JWT verification, Whitelist & Admin guard
│   │   └── scraper/              # RAWG, IGDB (Twitch OAuth) & Composite Scraper
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── components/           # Navbar, GameCard, GameDetailModal, AddGameModal
│   │   ├── composables/          # useAuth (GIS SDK), useGames (Download URL generator)
│   │   ├── types/                # TypeScript models
│   │   ├── App.vue
│   │   └── style.css             # Tailwind v4 import & custom styles
│   ├── vite.config.ts            # Tailwind v4 plugin + API proxy
│   └── package.json
├── Dockerfile                    # Multi-stage build (Node -> Go Static -> Alpine Runner)
├── docker-compose.yml            # Khởi chạy dịch vụ duy nhất + Volume /data
├── .env.example                  # Template biến môi trường
└── Makefile                      # Công cụ build & dev nhanh
```

---

## ⚙️ Cấu Hình Môi Trường (.env)

Sao chép `.env.example` thành `.env`:

```bash
cp .env.example .env
```

| Biến Môi Trường | Mô Tả | Mặc Định / Ví Dụ |
| :--- | :--- | :--- |
| `PORT` | Cổng HTTP Web App lắng nghe | `8080` |
| `DB_PATH` | Đường dẫn lưu SQLite database | `/data/catalog.db` (Docker) hoặc `data/catalog.db` (Local) |
| `GOOGLE_CLIENT_ID` | OAuth 2.0 Client ID tạo từ Google Cloud Console | `xxxx.apps.googleusercontent.com` |
| `WHITELIST_EMAILS` | Danh sách Gmail được phép truy cập | `alice@gmail.com,bob@gmail.com` |
| `ADMIN_EMAILS` | Danh sách Gmail có quyền thêm/xóa game | `alice@gmail.com` |
| `JWT_SECRET` | Khóa bí mật ký JWT Session | Tự động tạo random nếu để trống |
| `RAWG_API_KEY` | API Key miễn phí từ [rawg.io/apidocs](https://rawg.io/apidocs) | `abcdef...` |
| `TWITCH_CLIENT_ID` | Client ID cho IGDB fallback (tùy chọn) | `...` |
| `TWITCH_CLIENT_SECRET`| Client Secret cho IGDB (tùy chọn) | `...` |

> **Lưu ý cấu hình Google Cloud Console:**
> 1. Truy cập [Google Cloud Console Credentials](https://console.cloud.google.com/apis/credentials).
> 2. Tạo một **OAuth 2.0 Client ID** loại **Web application**.
> 3. Trong mục **Authorized JavaScript origins**, thêm địa chỉ domain/IP của bạn (ví dụ: `http://localhost:8080` hoặc `https://kame.yourdomain.com`).
> 4. Trong Google Drive, cấp quyền **Viewer** (Người xem) thủ công cho các email người dùng trong thư mục lưu game.

---

## 🚀 Triển Khai Docker (Single Container)

### 1. Build & Chạy qua Docker Compose

```bash
# Build image và khởi chạy ngầm
docker compose up -d --build

# Xem logs hoạt động
docker compose logs -f
```

### 2. Chạy Trực Tiếp Qua Docker CLI

```bash
# Build Docker image
docker build -t kame:latest .

# Run container với volume mount cho SQLite
docker run -d \
  --name kame_catalog \
  --restart unless-stopped \
  -p 8080:8080 \
  --env-file .env \
  -v $(pwd)/data:/data \
  kame:latest
```

Hệ thống sẽ sẵn sàng tại: **`http://localhost:8080`**

---

## 🛠️ Phát Triển Cục Bộ (Local Development)

Nếu bạn muốn chạy thử và debug trên máy local (Arch Linux / macOS / Linux):

```bash
# 1. Chạy Backend (Terminal 1)
make dev-backend

# 2. Chạy Frontend với Vite Hot-Reload (Terminal 2)
make dev-frontend
```

Frontend Vite sẽ chạy tại `http://localhost:5173` và tự động proxy các request `/api/*` về Go backend ở port `8080`.

Để build một file thực thi duy nhất cục bộ:

```bash
make build
# Binary tạo ra tại ./bin/kame (~13MB)
./bin/kame
```
