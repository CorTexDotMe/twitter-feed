FROM golang:1.23.0-alpine3.20 as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/main.go

FROM scratch

COPY --from=build /server /server

ENTRYPOINT ["/server"]