# MemKV - A Redis-like Key-Value Store in Go

This is a simple implementation of a Redis-like key-value store in Go. It provides basic functionality similar to Redis, including:

- String operations (GET, SET)
- List operations (LPUSH, RPUSH, LRANGE)
- Hash operations (HSET, HGET)
- Basic persistence

## Features

- In-memory key-value storage
- TCP server implementation
- RESP (Redis Serialization Protocol) support
- Basic data types support (Strings, Lists, Hashes)

## Getting Started

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Run `go run main.go` to start the server
4. Connect using any Redis client (default port: 6379)

## Project Structure

- `main.go`: Main server implementation
- `server/`: Server-related code
- `storage/`: Storage engine implementation
- `protocol/`: RESP protocol implementation

## License

MIT License 
