FROM golang:1.23.0-alpine3.20 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /message-server ./cmd/message-server/main.go

FROM scratch

COPY --from=build /message-server /message-server

ENTRYPOINT ["/message-server"]