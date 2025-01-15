# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /app

# dependencies cache
COPY go.mod go.sum ./
RUN go mod download && go mod verify && \
    rm -rf /root/.cache/go-build /go/pkg/mod


COPY cmd ./cmd
COPY internal ./internal

# generate HTMX files
RUN go install github.com/a-h/templ/cmd/templ@latest && \
    templ generate

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Development stage
FROM golang:1.23-alpine AS dev

WORKDIR /app

RUN go install github.com/air-verse/air@latest

# Copy the built binary and required directories
COPY --from=build /app/main /app/main
COPY cmd ./cmd
COPY internal ./internal

CMD ["air"]

# Production stage
FROM gcr.io/distroless/base:nonroot AS prod

WORKDIR /app

# Copy only the binary from the build stage
COPY --from=build /app/main /app/main

# Environment configuration
ENV PORT=8080
EXPOSE 8080

CMD ["./main"]
