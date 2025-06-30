# --- Stage 1: Build a aplicação Go ---
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# Copia os arquivos go.mod e go.sum para baixar as dependências primeiro
# Isso aproveita o cache do Docker se as dependências não mudarem
COPY go.mod go.sum ./
RUN go mod download

# Copia o restante do código fonte da sua API
COPY . .

# Constrói o executável da sua aplicação
# CGO_ENABLED=0 desabilita o cgo, tornando o binário estaticamente linkado
# -o api_server define o nome do executável final
# ./cmd/api é um exemplo do caminho para seu arquivo main.go. Ajuste conforme necessário.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o api_server ./cmd/

# --- Stage 2: Cria a imagem final leve ---
FROM alpine:latest

WORKDIR /app

# Instala os pacotes de dados de fuso horário (timezone) e certificados CA.
# tzdata é necessário para resolver "America/Sao_Paulo".
RUN apk --no-cache add tzdata ca-certificates

# Copia o executável construído da etapa anterior
COPY --from=builder /app/api_server .

# Copia quaisquer outros arquivos necessários pela sua API (ex: templates, arquivos de configuração estáticos)
# Exemplo: COPY --from=builder /app/configs ./configs

# Expõe a porta que sua API Gin escuta (padrão geralmente é 8080)
EXPOSE 8080

# Comando para rodar o executável da sua API
CMD ["./api_server"]