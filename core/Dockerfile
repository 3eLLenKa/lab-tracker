FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/lab-tracker ./cmd/app


FROM debian:bookworm-slim

RUN apt-get update \
	&& apt-get install -y --no-install-recommends ca-certificates \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /out/lab-tracker /app/lab-tracker
COPY config/config.yaml /app/config/config.yaml

RUN chown -R nobody:nogroup /app

ENV CONFIG_PATH=/app/config/config.yaml

USER nobody:nogroup

ENTRYPOINT ["/app/lab-tracker"]