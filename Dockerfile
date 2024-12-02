# Etapa de construcción: usa la imagen de Go para compilar la aplicación
FROM golang:1.23.0 AS builder

WORKDIR /app

# Copiar go.mod y go.sum primero para aprovechar la caché
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente al contenedor
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o alumnos-back ./cmd/main.go

# Etapa final: Crear imagen mínima con Alpine
FROM alpine:latest

WORKDIR /app

# Instalar bash (y otras utilidades si es necesario)
RUN apk add --no-cache bash

# Copiar el binario y el script wait-for-it.sh
COPY --from=builder /app/alumnos-back /app/alumnos-back
COPY wait-for-it.sh /app/wait-for-it.sh

# Establecer permisos ejecutables para wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

# Exponer el puerto en el que corre tu aplicación
EXPOSE 8080

# Usar wait-for-it.sh para esperar la base de datos antes de iniciar la app
ENTRYPOINT ["/app/wait-for-it.sh", "db:5432", "--", "/app/alumnos-back"]
