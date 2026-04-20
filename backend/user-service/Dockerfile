FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/server .

COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["./server"]
