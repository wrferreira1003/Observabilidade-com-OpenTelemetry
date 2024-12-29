# Etapa de construção
FROM golang:1.23 as builder

WORKDIR /app

# Copiar os arquivos do projeto
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Construir o binário para Linux ARM64
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o cloudrun ./cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun ./cmd/main.go

# Etapa final para execução
FROM alpine:latest
WORKDIR /app

# Instalar certificados CA (necessário para HTTPS)
RUN apk add --no-cache ca-certificates

# Copiar o binário e outros arquivos necessários
COPY --from=builder /app/cloudrun .
COPY --from=builder /app/.env .

ENTRYPOINT ["./cloudrun"]






