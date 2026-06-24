FROM node:24-alpine AS ui-builder

WORKDIR /src/ui

RUN corepack enable

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY ui/ ./
RUN pnpm run build

FROM golang:1.25-alpine AS server-builder

WORKDIR /src

COPY types/ types/
COPY server/ server/

RUN cd server && GOWORK=off go mod download
RUN cd server && GOWORK=off CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/aiusage-server .
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@v3.27.1

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && adduser -D -H -u 10001 aiusage

COPY --from=server-builder /out/aiusage-server /usr/local/bin/aiusage-server
COPY --from=server-builder /go/bin/goose /usr/local/bin/goose
COPY --from=ui-builder /src/ui/dist /app/ui/dist
COPY server/migrations /app/migrations

USER aiusage

ENV ENV=production \
    PORT=8080 \
    STATIC_DIR=/app/ui/dist

EXPOSE 8080

CMD ["aiusage-server"]
