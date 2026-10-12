.PHONY: all build dev-frontend dev-backend docker-build docker-up docker-down docker-logs clean

all: build

# Build Frontend and embed static assets into Go binary
build:
	@echo "==> Building Frontend (Vite + Tailwind v4)..."
	cd frontend && pnpm build
	@echo "==> Syncing static assets to Backend embed directory..."
	rm -rf backend/cmd/server/dist && mkdir -p backend/cmd/server/dist
	cp -r frontend/dist/* backend/cmd/server/dist/
	@echo "==> Building Go Binary (kame)..."
	cd backend && go build -ldflags="-s -w" -o ../bin/kame cmd/server/main.go
	@echo "==> Build complete! Output binary: ./bin/kame"

# Start Frontend development server (port 5173 with proxy to :8080)
dev-frontend:
	cd frontend && pnpm dev

# Start Backend development server (port 8080)
dev-backend:
	cd backend && PORT=8080 DB_PATH=data/catalog.db go run cmd/server/main.go

# Docker Single Image build & container management
docker-build:
	docker build -t kame:latest .

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

clean:
	rm -rf bin/ frontend/dist/ backend/cmd/server/dist/ data/
