FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN GOARCH=amd64 go install github.com/pressly/goose/v3/cmd/goose@v3.27.1

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o jarvis-api ./cmd/api

FROM alpine:3.21 AS runner

WORKDIR /app

RUN addgroup -S jarvis && adduser -S jarvis -G jarvis

COPY --from=builder /app/jarvis-api .

COPY --from=builder /app/migrations ./migrations

COPY --from=builder /go/bin/goose .

COPY entrypoint.sh .
RUN chown jarvis:jarvis entrypoint.sh jarvis-api goose
RUN chmod +x entrypoint.sh jarvis-api goose

USER jarvis

EXPOSE 8080

CMD ["./entrypoint.sh"]
