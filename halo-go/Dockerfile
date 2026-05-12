FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o halo-server cmd/server/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata sqlite-libs && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=builder /app/halo-server /app/halo-server

ENV HALO_PORT=8090 \
    HALO_DB_DATASOURCE=/data/halo.db \
    HALO_WORK_DIR=/data

EXPOSE 8090

VOLUME ["/data"]

ENTRYPOINT ["/app/halo-server"]
