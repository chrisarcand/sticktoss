# Multi-stage Dockerfile for Stick Toss

# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm install

COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# Install build dependencies for CGO (required for SQLite)
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy backend source
COPY backend/ ./

# Tidy dependencies to generate go.sum
RUN go mod tidy

# Download dependencies
RUN go mod download

# Build the Go binary
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/server

# Stage 3: Final runtime image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS and SQLite support
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy the Go binary from backend-builder
COPY --from=backend-builder /app/server /app/server

# Copy the frontend build from frontend-builder
COPY --from=frontend-builder /app/frontend/dist /app/frontend/dist

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Run the server
CMD ["/app/server"]
