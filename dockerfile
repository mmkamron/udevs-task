FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -ldflags='-s -w' -o main ./cmd/api/

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .

CMD ["./main"]
