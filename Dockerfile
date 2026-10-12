# ==========================================
# STAGE 1: Build Frontend (Vue 3 + Vite + Tailwind v4)
# ==========================================
FROM node:22-alpine AS frontend-builder
WORKDIR /app/frontend

RUN npm install -g pnpm

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY frontend ./
RUN pnpm build

# ==========================================
# STAGE 2: Build Backend (Go Static Binary with Embedded SPA)
# ==========================================
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app/backend

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend ./

# Copy compiled frontend dist into embed directory
RUN rm -rf cmd/server/dist && mkdir -p cmd/server/dist
COPY --from=frontend-builder /app/frontend/dist ./cmd/server/dist

# Build pure Go static binary (No CGO required for modernc.org/sqlite)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -o /app/kame-server \
    ./cmd/server/main.go

# ==========================================
# STAGE 3: Production Runner (Minimal Single Container)
# ==========================================
FROM alpine:3.21 AS runner

RUN apk add --no-cache ca-certificates tzdata dumb-init sqlite

# Create unprivileged user & data mount directory
RUN addgroup -S appgroup && adduser -S appuser -G appgroup && \
    mkdir -p /data && chown -R appuser:appgroup /data

WORKDIR /app

COPY --from=backend-builder --chown=appuser:appgroup /app/kame-server ./kame-server

USER appuser

ENV PORT=8080
ENV DB_PATH=/data/catalog.db

VOLUME ["/data"]
EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD ["./kame-server"]
