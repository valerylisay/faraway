FROM golang:1.24 AS builder

WORKDIR /app

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/server main.go


FROM alpine:latest

ENV ADDR=:8080 DIFFICULTY=4

COPY --from=builder /bin/server /bin/server

EXPOSE 8080
CMD ["sh", "-c", "/bin/server --addr=${ADDR} --difficulty=${DIFFICULTY}"]
