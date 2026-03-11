# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend (Pure Go, no CGO needed)
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .

# Stage 3: Production
FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata curl
WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/backend/server .
RUN chmod +x ./server

# Copy frontend dist
COPY --from=frontend-builder /app/frontend/dist ./static

# Create data directory for SQLite
RUN mkdir -p /app/data

# Environment variables
ENV PORT=8080
ENV DB_PATH=/app/data/cvbuilder.db
ENV GIN_MODE=release
ENV TZ=Asia/Baghdad

# Health check (use shell form so $PORT is expanded at runtime)
HEALTHCHECK --interval=10s --timeout=5s --start-period=15s --retries=5 \
  CMD curl -f http://127.0.0.1:${PORT:-8080}/api/v1/stats || exit 1

EXPOSE 8080

CMD ["./server"]
