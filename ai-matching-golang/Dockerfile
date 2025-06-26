FROM golang:1.24-alpine AS builder

# git と make をインストール
RUN apk add --no-cache git make curl

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin/migrate

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY . .

# Generate sqlc code
RUN sqlc generate

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest AS final

# 実行環境にも make を配置
RUN apk --no-cache add ca-certificates make

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env* ./
COPY --from=builder /app/db/migrations ./db/migrations
COPY --from=builder /app/Makefile ./

COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

EXPOSE 8080

CMD ["./main"]