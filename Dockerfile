FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /calc-service ./cmd/calc_service/

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /calc-service /calc-service
COPY migrations/init.sql /init.sql

EXPOSE 8080

CMD ["./calc-service"]