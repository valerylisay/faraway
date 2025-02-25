FROM golang:1.24 AS builder

WORKDIR /app

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/client main.go


FROM alpine:latest

ENV SERVER_ADDR=:8080

COPY --from=builder /bin/client /bin/client

CMD ["sh", "-c", "/bin/client --server-addr=${SERVER_ADDR}"]
