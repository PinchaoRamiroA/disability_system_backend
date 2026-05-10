FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main ./cmd/api

CMD ["./main"]

