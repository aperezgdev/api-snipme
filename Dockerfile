FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/main ./src/cmd/main.go

FROM alpine:3.23

WORKDIR /app
COPY db/geo/GeoLite2-Country.mmdb db/geo/GeoLite2-Country.mmdb
COPY --from=builder /app/main /snipme-api

RUN adduser -D snipmeuser

USER snipmeuser

ENTRYPOINT [ "/snipme-api" ]
