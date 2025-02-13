# Этап сборки (build stage)
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /short-links ./cmd/server

# Финальный этап (final stage)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /short-links .
COPY .env.example .env

EXPOSE 8080
CMD ["./short-links", "-storage", "postgres"]