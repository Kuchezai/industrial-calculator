FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o industrial-calculator ./cmd/main.go
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/industrial-calculator .
EXPOSE 8080 50051
CMD ["./industrial-calculator"]