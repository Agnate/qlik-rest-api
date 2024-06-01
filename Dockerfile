FROM golang:1.22-alpine as base

# DEV environment
FROM base as dev
WORKDIR /app
# Use air to watch for file reloading during development
RUN go install github.com/cosmtrek/air@latest

# TODO: Change this to only use go.mod and go.sum
#COPY . /app/
COPY go.mod go.sum .env ./
RUN go mod download

# Other environments
# Do we want to build swagger API separately?
#RUN go install github.com/swaggo/swag/cmd/swag@latest