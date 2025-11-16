FROM golang:1.24.4 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o test ./cmd

FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /app/test .
COPY --from=builder /app/config.yaml ./


EXPOSE 8080

CMD ["./test"]