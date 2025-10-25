FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./src/cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/main /snipme-api

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

RUN useradd -m snipmeuser

USER snipmeuser

ENTRYPOINT [ "/snipme-api" ]
