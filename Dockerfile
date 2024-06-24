FROM golang:1.22-alpine

WORKDIR /app

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Instalar bash
RUN apk update && apk add --no-cache bash

# Instalar air
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

RUN go get -u golang.org/x/crypto/bcrypt
RUN go get -u github.com/gofiber/swagger
RUN go get -u github.com/gofiber/template@v1.8.3
RUN go get -u github.com/golang-jwt/jwt/v5

# Copiar el resto del código
COPY . .

# Descargar dependencias de Go
RUN go mod tidy

# Comando para ejecutar air y tu aplicación
# CMD ["air", "-c", ".air.toml"]
