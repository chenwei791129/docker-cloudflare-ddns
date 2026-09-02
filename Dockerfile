FROM golang:1.27-alpine AS builder

WORKDIR /build

COPY go.mod ./
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o cloudflare-ddns .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/cloudflare-ddns /cloudflare-ddns

ENTRYPOINT ["/cloudflare-ddns"]
