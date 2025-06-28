# Triff

[![Go Reference](https://pkg.go.dev/badge/github.com/nitrix4ly/triff.svg)](https://pkg.go.dev/github.com/nitrix4ly/triff)
[![Build Status](https://github.com/nitrix4ly/triff/actions/workflows/go.yml/badge.svg)](https://github.com/nitrix4ly/triff/actions)
[![License](https://img.shields.io/github/LICENSE/nitrix4ly/triff)](LICENSE)

Triff is a lightweight in-memory key-value database engine written in Go. It provides support for multiple data types including strings, sets, hashes, and lists, along with optional persistence using a disk-based engine. Triff is modular, testable, and designed to be extensible and educational.

## Features

- In-memory key-value store with fast access
- Support for data types: string, set, list, hash
- Modular and clean architecture
- Optional persistence using JSON disk engine
- Example implementations including HTTP and Discord bot integration
- Fully tested and designed for learning or extension

## Installation

To use Triff as a Go module, run:

```bash
go get github.com/nitrix4ly/triff
```

Or add it to your `go.mod` file:

```go
require github.com/nitrix4ly/triff v0.0.0
```

## Basic Usage

```go
import (
    "fmt"
    "github.com/nitrix4ly/triff/core"
)

func main() {
    db := core.NewDatabase(&core.Config{})
    value := &core.TriffValue{
        Type: core.STRING,
        Data: "Example",
    }
    db.Set("exampleKey", value)

    result, ok := db.Get("exampleKey")
    if ok {
        fmt.Println(result.Data)
    }
}
```

## Module Structure

- `core/` - Main database and configuration interfaces
- `commands/` - Built-in operations for supported data types
- `storage/` - In-memory and disk storage engines
- `server/` - HTTP and TCP server implementations
- `utils/` - Utilities for parsing and configuration
- `examples/` - Sample applications and usage patterns

## Testing

To run the test suite:

```bash
go test ./...
```

## Building

To build the project:

```bash
go build ./...
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
