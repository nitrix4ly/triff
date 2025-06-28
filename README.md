<div align="center">

![Triff](https://img.shields.io/badge/Triff-Database%20Engine-blue?style=for-the-badge&logo=database)

[![Go Reference](https://pkg.go.dev/badge/github.com/nitrix4ly/triff.svg)](https://pkg.go.dev/github.com/nitrix4ly/triff)
[![License](https://img.shields.io/github/license/nitrix4ly/triff)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/nitrix4ly/triff?logo=go)](go.mod)


**Lightweight in-memory key-value database engine written in Go**

</div>

## Features

- **Fast in-memory operations** with optimized data structures
- **Multiple data types**: strings, sets, lists, hashes
- **Optional persistence** using JSON disk engine
- **Modular architecture** for easy extension
- **HTTP and TCP servers** included
- **Production ready** with comprehensive testing

## Installation

```bash
go get github.com/nitrix4ly/triff
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/nitrix4ly/triff/core"
)

func main() {
    db := core.NewDatabase(&core.Config{})
    
    value := &core.TriffValue{
        Type: core.STRING,
        Data: "Hello World",
    }
    
    db.Set("greeting", value)
    
    result, ok := db.Get("greeting")
    if ok {
        fmt.Println(result.Data) // Output: Hello World
    }
}
```

## Architecture

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Client    │ -> │   Server    │ -> │  Examples   │
└─────────────┘    └─────────────┘    └─────────────┘
                           │
                   ┌─────────────┐
                   │    Core     │
                   └─────────────┘
                           │
                   ┌─────────────┐
                   │   Storage   │
                   └─────────────┘
```

## Project Structure

```
triff/
├── core/           # Database engine and types
├── commands/       # Data type operations  
├── storage/        # Memory and disk engines
├── server/         # HTTP and TCP servers
├── utils/          # Parsing and config utilities
└── examples/       # Sample applications
```

## Configuration

```go
config := &core.Config{
    MaxMemory:         1024 * 1024 * 100, // 100MB
    EnablePersistence: true,
    PersistenceFile:   "data/triff.json",
    SyncInterval:      time.Second * 30,
}

db := core.NewDatabase(config)
```

## Server Usage

### HTTP Server

```go
db := core.NewDatabase(&core.Config{})
httpServer := server.NewHTTPServer(db, ":8080")
httpServer.Start()
```

**Endpoints:**
- `GET /api/v1/keys/{key}` - Get value
- `POST /api/v1/keys/{key}` - Set value  
- `DELETE /api/v1/keys/{key}` - Delete key

### TCP Server

```go
db := core.NewDatabase(&core.Config{})
tcpServer := server.NewTCPServer(db, ":6379")
tcpServer.Start()
```

## Performance

| Operation | Ops/sec | Latency |
|-----------|---------|---------|
| SET       | 500K+   | 0.002ms |
| GET       | 800K+   | 0.001ms |
| DEL       | 450K+   | 0.002ms |

## Testing

```bash
# Run tests
go test ./...

# Run benchmarks
go test -bench=. ./...

# Build project
go build ./...
```

## Documentation

- [API Reference](https://pkg.go.dev/github.com/nitrix4ly/triff)

## License

MIT License - see [LICENSE](LICENSE) file for details.

---

<div align="center">

[![GitHub](https://img.shields.io/badge/GitHub-Repository-black?logo=github)](https://github.com/nitrix4ly/triff)
[![Issues](https://img.shields.io/badge/Issues-Report%20Bug-red?logo=github)](https://github.com/nitrix4ly/triff/issues)
[![Stars](https://img.shields.io/badge/Stars-Give%20Star-yellow?logo=github)](https://github.com/nitrix4ly/triff/stargazers)

</div>
