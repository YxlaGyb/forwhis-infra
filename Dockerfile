FROM golang:1.26.0-alpine AS builder

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download

RUN go build -o forestwhisperGo

FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /app

COPY --from=builder /app/forestwhisperGo /app/forestwhisperGo

EXPOSE 2324

CMD ["/app/forestwhisperGo"]
