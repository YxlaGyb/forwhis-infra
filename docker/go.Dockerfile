FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY apps ./apps
RUN mkdir -p /out && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/forestwhisperGo ./apps/cmd/servers

FROM alpine:3.23

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /app \
    && chown app:app /app

WORKDIR /app
COPY --from=builder /out/forestwhisperGo /usr/local/bin/forestwhisperGo

USER app
EXPOSE 2324
CMD ["/usr/local/bin/forestwhisperGo"]
