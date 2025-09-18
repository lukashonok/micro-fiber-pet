# Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем модули
COPY go.mod go.sum ./
RUN go mod tidy

# Копируем весь проект
COPY . .

# Билдим main.go из cmd
RUN go build -o app ./cmd

# Финальный минимальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 3000
CMD ["./app"]