FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /worker ./internal/agent/worker/

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /worker /worker
CMD ["./worker"]