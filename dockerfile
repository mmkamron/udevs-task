FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN go build -ldflags='-s -w' -o main ./cmd/api/

CMD ["./main"]
