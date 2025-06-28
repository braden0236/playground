FROM --platform=$BUILDPLATFORM golang:1.24.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /go/bin/go-grpc ./cmd/go-grpc

FROM alpine:3.17
WORKDIR /app

COPY --from=builder /go/bin/go-grpc ./

EXPOSE 8080
ENTRYPOINT ["/app/go-grpc"]
