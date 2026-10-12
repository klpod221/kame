.PHONY: all build dev-frontend dev-backend docker-build docker-up docker-down clean

all: build

# Build Frontend rồi nhúng vào Go binary
build:
	@echo "==> Building Frontend (Vite + Tailwind v4)..."
	cd frontend && pnpm build
	@echo "==> Syncing static assets to Backend embed..."
	rm -rf backend/cmd/server/dist && mkdir -p backend/cmd/server/dist
	cp -r frontend/dist/* backend/cmd/server/dist/
	@echo "==> Building Go Binary (kame)..."
	cd backend && go build -ldflags="-s -w" -o ../bin/kame cmd/server/main.go
	@echo "==> Build complete! Output binary: ./bin/kame"

# Chạy frontend dev server (port 5173, proxy sang :8080)
dev-frontend:
	cd frontend && pnpm dev

# Chạy backend dev server (port 8080)
dev-backend:
	cd backend && PORT=8080 DB_PATH=data/catalog.db go run cmd/server/main.go

# Docker Single Image build & run
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
