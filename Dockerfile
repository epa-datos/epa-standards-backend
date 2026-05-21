# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o api ./cmd/api

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/api .
EXPOSE 8080
HEALTHCHECK --interval=30s CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./api"]
