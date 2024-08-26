FROM golang:1.23.0-alpine3.20 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /message-consumer ./cmd/message-consumer/main.go

FROM scratch

COPY --from=build /message-consumer /message-consumer

ENTRYPOINT ["/message-consumer"]