ARG GO_VERSION=1.26

FROM oven/bun:1 AS ui-builder

WORKDIR /ui

COPY ui/package.json ui/bun.lock* ./
RUN bun install --frozen-lockfile

COPY ui/ ./
RUN bun run build

FROM golang:${GO_VERSION}-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=ui-builder /ui/dist ./ui/dist
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/substrack .

FROM debian:trixie-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
RUN useradd -m -u 10001 appuser && mkdir -p /app/data && chown -R appuser:appuser /app

COPY --from=builder /app/substrack /app/substrack

EXPOSE 8080
VOLUME ["/app/data"]
USER appuser

CMD ["/app/substrack"]
