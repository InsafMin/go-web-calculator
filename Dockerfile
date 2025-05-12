FROM golang:1.23 AS builder

RUN apt-get update && apt-get install -y gcc musl-dev

WORKDIR /app
COPY . .

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

ENV CGO_ENABLED=1
RUN go build -ldflags="-extldflags=-static" -o /calc_service ./cmd/calc_service
RUN go build -ldflags="-extldflags=-static" -o /worker ./cmd/worker

FROM alpine:3.17
WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=builder /calc_service /app/
COPY --from=builder /worker /app/
COPY migrations /app/migrations

EXPOSE 50051
EXPOSE 8080