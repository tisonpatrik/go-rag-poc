# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify && \
    rm -rf /root/.cache/go-build /go/pkg/mod

# Copy application source
COPY cmd ./cmd
COPY internal ./internal

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/main ./cmd/api/main.go

# Development stage
FROM golang:1.23-alpine AS dev

WORKDIR /app

# Install air for live reloading (use explicit version)
RUN go install github.com/air-verse/air@latest

RUN air -v

# Copy the source and binary
COPY --from=build /app/build/main /app/build/main
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY .air.toml ./

CMD ["air", "-c", ".air.toml"]

# Production stage
FROM gcr.io/distroless/base:nonroot AS prod

WORKDIR /app

# Copy the built binary
COPY --from=build /app/build/main /app/main

# Environment configuration
ENV PORT=8080
EXPOSE 8080

CMD ["/app/main"]
