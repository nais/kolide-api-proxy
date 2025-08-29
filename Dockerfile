ARG GO_VERSION="1.25"
FROM golang:${GO_VERSION} AS builder
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o ./bin/proxy main.go

FROM gcr.io/distroless/base
WORKDIR /app
COPY --from=builder /src/bin/proxy /app/proxy
CMD ["/app/proxy"]