# --- Stage 1: Build a aplicação Go ---
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o api_server ./cmd/

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add tzdata ca-certificates

COPY --from=builder /app/api_server .

EXPOSE 8080

# Comando para rodar o executável da sua API
CMD ["./api_server"]