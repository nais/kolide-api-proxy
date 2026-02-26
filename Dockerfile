ARG GO_VERSION="1.26"
FROM golang:${GO_VERSION} AS builder
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o ./bin/proxy main.go

FROM scratch
ADD https://curl.haxx.se/ca/cacert.pem /etc/ssl/certs/ca-certificates.crt
WORKDIR /app
COPY --from=builder /src/bin/proxy /app/proxy
CMD ["/app/proxy"]
