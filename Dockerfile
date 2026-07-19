FROM golang:bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o tandoor-mcp ./src/

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates netcat-openbsd \
	&& rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/tandoor-mcp .

EXPOSE 8080
ENTRYPOINT ["./tandoor-mcp"]
CMD ["-transport", "sse", "-port", "8080", "-host", "0.0.0.0"]
