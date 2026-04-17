FROM node:22-alpine AS web-builder

WORKDIR /app/web
COPY web/package.json web/package-lock.json* ./
RUN if [ -f package-lock.json ]; then npm ci; else npm install; fi
COPY web/ ./
RUN npm run build

FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /halo ./cmd/halo

FROM alpine:3.22

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

ENV HALO_WORK_DIR=/root/.halo2 \
    HALO_ADDR=:8090 \
    HALO_DB_PATH=/root/.halo2/db/halo.db \
    TZ=Asia/Shanghai

COPY --from=builder /halo /usr/local/bin/halo
COPY --from=web-builder /app/web/dist ./web/dist

EXPOSE 8090

ENTRYPOINT ["/usr/local/bin/halo"]
