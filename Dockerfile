# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.25-alpine AS backend-builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o server .

# Stage 3: Production
FROM alpine:3.20
RUN apk add --no-cache ca-certificates sqlite-libs tzdata curl
WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/backend/server .

# Copy frontend dist
COPY --from=frontend-builder /app/frontend/dist ./static

# Create data directory for SQLite
RUN mkdir -p /app/data

# Environment variables
ENV PORT=8080
ENV DB_PATH=/app/data/cvbuilder.db
ENV GIN_MODE=release
ENV TZ=Asia/Baghdad

# Health check
HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD curl -f http://localhost:8080/api/v1/stats || exit 1

EXPOSE 8080

CMD ["./server"]
