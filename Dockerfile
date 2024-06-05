FROM golang:1.22-alpine as base

# DEV environment
FROM base as dev
WORKDIR /api

# Use air to watch for file reloading during development
RUN go install github.com/air-verse/air@latest

# Copy modules and environment variables, then install mods
COPY go.mod go.sum .env ./
COPY migrations/ ./migrations/
RUN go mod download && go mod verify