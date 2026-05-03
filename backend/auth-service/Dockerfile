FROM golang:1.25-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    pkg-config \
    gcc \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go
RUN go build -o worker ./cmd/worker/main.go

FROM debian:bookworm-slim

WORKDIR /root

RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server .
COPY --from=builder /app/worker .

COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["./server"]
