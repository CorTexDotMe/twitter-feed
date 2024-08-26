FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /cmd/message-consumer ./cmd/message-consumer/main.go

CMD ["/cmd/message-consumer"]