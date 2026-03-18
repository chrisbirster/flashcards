FROM oven/bun:1.3.6 AS web-builder
WORKDIR /app/web
COPY web/package.json web/bun.lock* ./
COPY web/tsconfig*.json ./
COPY web/vite.config.ts ./
COPY web/eslint.config.js ./
COPY web/index.html ./
COPY web/public ./public
COPY web/src ./src
COPY web/test ./test
COPY web/scripts ./scripts
RUN bun install
RUN bun run build

FROM golang:1.25-bookworm AS go-builder
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends build-essential ca-certificates && rm -rf /var/lib/apt/lists/*
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=web-builder /app/web/dist ./web/dist
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /out/vutadex .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=go-builder /out/vutadex /app/vutadex
ENV PORT=8080
EXPOSE 8080
CMD ["/app/vutadex"]
