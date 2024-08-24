FROM golang:1.23.0-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /cmd/twitterfeed/main ./cmd/twitterfeed/main.go

CMD ["/cmd/twitterfeed/main"]