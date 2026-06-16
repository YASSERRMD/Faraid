# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /faraid ./cmd/faraidd

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
RUN addgroup -S faraid && adduser -S faraid -G faraid
USER faraid
COPY --from=builder /faraid /usr/local/bin/faraid
EXPOSE 8080
ENTRYPOINT ["faraid"]
