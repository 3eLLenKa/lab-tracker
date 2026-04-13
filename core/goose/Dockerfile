FROM golang:1.25-alpine

RUN apk add --no-cache git

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY migrations /app/migrations

ENTRYPOINT ["goose"]
