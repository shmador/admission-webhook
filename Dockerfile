# 1) Build stage: compile static binary
FROM golang:1.24 AS builder
WORKDIR /app

# copy go.mod/go.sum first to cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of your code
COPY . .

# disable CGO so no libc is linked, and strip symbols
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o webhook

# 2) Final stage: scratch (empty) image
FROM scratch
# expose port for webhook server
EXPOSE 8443
# copy the static binary
COPY --from=builder /app/webhook /webhook

ENTRYPOINT ["/webhook"]

