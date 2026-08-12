# --- Build stage ---------------------------------------------------------
FROM golang:1.23-alpine AS build
WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/epa-standards-backend .

# --- Runtime stage ---------------------------------------------------------
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=build /app/bin/epa-standards-backend .

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Real config comes from env vars injected by the platform (Cloud Run,
# docker run --env-file, k8s secrets, ...). No .env file is copied here.
ENTRYPOINT ["/app/epa-standards-backend"]
