# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copying sources and build them
COPY . .
RUN go build -o app

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
EXPOSE 3000
CMD ["./app"]