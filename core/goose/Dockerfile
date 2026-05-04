FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates wget

RUN wget -O /usr/local/bin/goose \
    https://github.com/pressly/goose/releases/download/v3.26.0/goose_linux_x86_64 \
    && chmod +x /usr/local/bin/goose

COPY migrations /app/migrations

ENTRYPOINT ["goose"]