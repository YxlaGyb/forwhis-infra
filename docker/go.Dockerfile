FROM docker.1ms.run/library/golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY apps/api ./apps/api
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    mkdir -p /out && CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" \
    -o /out/forwhis-infra ./apps/api/cmd/servers

FROM docker.1ms.run/library/alpine:3.23

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app \
    && mkdir -p /app \
    && chown app:app /app

WORKDIR /app
COPY --from=builder /out/forwhis-infra /usr/local/bin/forwhis-infra

USER app
EXPOSE 23324
CMD ["/usr/local/bin/forwhis-infra"]
