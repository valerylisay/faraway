# Word of Wisdom TCP server

Implementation of “Word of Wisdom” tcp server, protected from DDOS attacks with the Proof of Work challenge-response protocol. I used Hashcash algo because it's simple, well-known, effective and tunable.

## How to run

```bash
# Using docker compose
make up

# Then you can fetch quote one by one
make get_quote

# Standalone server (docker)
docker build -t wow-server:latest -f server/Dockerfile server
docker run --rm -p 8080:8080 -e ADDR=:8080 -e DIFFICULTY=4 wow-server:latest

# Standalone client (docker)
docker build -t wow-client:latest -f client/Dockerfile client
docker run --rm -e SERVER_ADDR=host.docker.internal:8080 wow-client:latest

# Standalone server
make build_server && server/server --addr=:8080 --difficulty=4

# Client
make build_client && client/client --server-addr=:8080

# Tests
make tests
```
