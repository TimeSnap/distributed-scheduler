# Use an official GoLang runtime as the base image
FROM golang:1.23-alpine AS builder

ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-build
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files to the working directory
COPY go.mod go.sum ./

# Download dependencies
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

# Verify dependencies
RUN go mod verify

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# install go swaggo, used to generate swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
# build docs
RUN swag init -g internal/api/http/doc.go

# Build the Go application
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o bin/runner cmd/runner/main.go

# Use an official Alpine Linux runtime as a base image
FROM alpine:latest

# Add curl for health checks
RUN apk --update --no-cache add curl

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the current stage
COPY --from=builder /app/bin/runner /app/runner

HEALTHCHECK --interval=5s --timeout=3s --retries=3  CMD curl --fail http://localhost:8000/healthz || exit 1

# set command to run when starting the container
CMD ["/app/runner"]