# -------- Build Stage --------
FROM golang:1.25.0-trixie AS builder

WORKDIR /app

# Install deps first (caching)
#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download

# Copy source
COPY . .

# Build statically linked binary
RUN go build -o server ./cmd/api

# -------- Run Stage --------
FROM debian:trixie-slim

# Install CA certificates (needed for outbound HTTPS calls)
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Default environment variables (can override at runtime)
ENV PORT=8080
ENV ENVIRONMENT=development

EXPOSE $PORT

ENTRYPOINT ./server --port=$PORT --env=$ENVIRONMENT