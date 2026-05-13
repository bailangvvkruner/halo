FROM golang:1.25-alpine3.21 AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o halo-server cmd/server/main.go

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata sqlite-libs curl && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

ENV HALO_PORT=8090 \
    HALO_DB_DATASOURCE=/data/halo.db \
    HALO_WORK_DIR=/data

WORKDIR /app

VOLUME ["/data"]

EXPOSE 8090

LABEL org.opencontainers.image.title="Halo Go" \
      org.opencontainers.image.description="Halo Blog CMS - Go 重写版" \
      org.opencontainers.image.vendor="halo-dev" \
      org.opencontainers.image.source="https://github.com/halo-dev/halo-go"

COPY --from=builder /app/halo-server /usr/local/bin/halo-server

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD curl -sf http://localhost:${HALO_PORT:-8090}/actuator/health || exit 1

ENTRYPOINT ["halo-server"]
CMD ["--halo.work-dir=/data", "--server.port=8090"]