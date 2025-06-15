# Dockerfile for dor-admission-webhook
FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o webhook .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/webhook .
ENTRYPOINT ["./webhook", "-tlsCertFile=/certs/tls.crt", "-tlsKeyFile=/certs/tls.key"]

