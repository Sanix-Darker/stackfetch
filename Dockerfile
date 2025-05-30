FROM golang:1.21-alpine AS builder

RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev

WORKDIR /app
COPY . .

# Build statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -a -o stackfetch ./cmd/stackfetch

FROM scratch
COPY --from=builder /app/stackfetch /stackfetch
ENTRYPOINT ["/stackfetch"]
