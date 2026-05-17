# Stage 1: Build Node.js frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /build/frontend

COPY frontend/package*.json ./

RUN npm ci

COPY frontend/ ./

RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.26-alpine AS backend-builder

WORKDIR /build/backend

COPY backend/go.mod backend/go.sum ./

RUN go mod download

COPY backend/cmd ./cmd
COPY backend/internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# Stage 3: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=backend-builder /build/backend/server .

COPY --from=frontend-builder /build/frontend/dist ./dist

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/ || exit 1

CMD ["./server"]
