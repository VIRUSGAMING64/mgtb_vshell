# syntax=docker/dockerfile:1

# --- build ---
FROM golang:1.26-bookworm AS build

WORKDIR /src

# Cache de dependencias
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Mismas flags que run.sh (sin el sudo, en el contenedor corremos como root)
RUN CGO_ENABLED=0 go build -tags netgo -ldflags '-s -w' -o /out/app .

# --- runtime ---
FROM debian:bookworm-slim

# ca-certificates: llamadas HTTPS a la API de Telegram
# ffmpeg: requerido por core/ffmpeg (compresión/split de video)
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates ffmpeg \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /out/app ./app
COPY web ./web
COPY admins.conf ./admins.conf

# Puerto que sirve web/serv.go
EXPOSE 80

CMD ["./app"]
