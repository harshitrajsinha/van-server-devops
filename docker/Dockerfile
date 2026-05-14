# Build stage
FROM golang:1.24.13-alpine3.23 AS builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

# Set env vars
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Build
RUN go build -o main handler.go

# Final stage
FROM alpine:3.21

RUN addgroup -S srvgrp && adduser -S srvusr -G srvgrp

WORKDIR /server

COPY --from=builder /app/main .
COPY --from=builder /app/store/schema.sql ./store/schema.sql

EXPOSE 8080

USER srvusr

ENTRYPOINT ["./main"]
