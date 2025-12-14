ARG GO_VERSION=1.25.5

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/substrack .

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
RUN useradd -m -u 10001 appuser && mkdir -p /app/data && chown -R appuser:appuser /app

COPY --from=builder /app/substrack /app/substrack

EXPOSE 8080
VOLUME ["/app/data"]
USER appuser

CMD ["/app/substrack"]
