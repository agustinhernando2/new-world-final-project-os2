FROM golang:1.22-alpine

WORKDIR /app

# Instalar dependencias del sistema
RUN apk add --no-cache git

# Instalar bash
RUN apk update && apk add --no-cache bash

# Instalar air
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código
COPY . .

# Descargar dependencias de Go
RUN go mod tidy

# Comando para ejecutar air y tu aplicación
CMD ["air", "-c", ".air.toml"]
