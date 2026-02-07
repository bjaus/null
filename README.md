# null

[![Go Reference](https://pkg.go.dev/badge/github.com/bjaus/null.svg)](https://pkg.go.dev/github.com/bjaus/null)
[![Go Report Card](https://goreportcard.com/badge/github.com/bjaus/null)](https://goreportcard.com/report/github.com/bjaus/null)
[![CI](https://github.com/bjaus/null/actions/workflows/ci.yml/badge.svg)](https://github.com/bjaus/null/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/bjaus/null/branch/main/graph/badge.svg)](https://codecov.io/gh/bjaus/null)

Generic nullable type for Go that distinguishes between unset, null, and valid values.

## Features

- **Three-State Logic** — Distinguish between "not provided", "explicitly null", and "has value"
- **Generic** — Single `Value[T]` type works with any type
- **JSON Support** — Full marshal/unmarshal with three-state preservation
- **SQL Support** — Implements `Scanner` and `Valuer` for all common types
- **DynamoDB Support** — Via `nullddb` subpackage
- **Zero Dependencies** — Core package uses only the standard library

## Installation

```bash
go get github.com/bjaus/null
```

For DynamoDB support:

```bash
go get github.com/bjaus/null/nullddb
```

## The Problem

Go can't distinguish between "field not provided" and "field explicitly set to null":

```go
type Request struct {
    Name *string `json:"name"`
}

// Both of these result in Name == nil:
// {"name": null}
// {}
```

This matters for PATCH APIs where you need to know:
- **Unset**: Client didn't mention it → don't change it
- **Null**: Client wants to clear it → set to NULL
- **Value**: Client wants to update it → set to value

## Quick Start

```go
package main

import (
    "encoding/json"
    "fmt"

    "github.com/bjaus/null"
)

type UpdateRequest struct {
    Name  null.Value[string] `json:"name"`
    Email null.Value[string] `json:"email"`
    Age   null.Value[int]    `json:"age"`
}

func main() {
    data := `{"name": "Alice", "email": null}`

    var req UpdateRequest
    json.Unmarshal([]byte(data), &req)

    fmt.Println(req.Name.IsValid())  // true - has value "Alice"
    fmt.Println(req.Email.IsNull())  // true - explicitly null
    fmt.Println(req.Age.IsSet())     // false - not in JSON
}
```

## Usage

### Creating Values

```go
// Valid value
name := null.New("Alice")

// Explicit null
name := null.NewNull[string]()

// From pointer (nil becomes null)
name := null.NewPtr(namePtr)

// Zero value is unset
var name null.Value[string]
```

### Checking State

```go
v.IsSet()   // true if null OR valid (was explicitly provided)
v.IsNull()  // true if explicitly null
v.IsValid() // true if has a real value
```

| State | IsSet | IsNull | IsValid |
|-------|-------|--------|---------|
| Unset | false | false  | false   |
| Null  | true  | true   | false   |
| Valid | true  | false  | true    |

### Extracting Values

```go
v.Get()          // Returns value or zero value of T
v.GetOr("default") // Returns value or the default
v.Ptr()          // Returns *T or nil
```

### PATCH Request Pattern

```go
func UpdateUser(req UpdateRequest, user *User) {
    if req.Name.IsSet() {
        if req.Name.IsNull() {
            user.Name = nil  // Clear the field
        } else {
            user.Name = req.Name.Ptr()  // Update the field
        }
    }
    // If !req.Name.IsSet(), leave unchanged
}
```

### SQL Integration

```go
type User struct {
    ID   int64
    Name null.Value[string]
}

// Scanning
row.Scan(&user.ID, &user.Name)
// NULL → user.Name.IsNull() == true
// "Alice" → user.Name.Get() == "Alice"

// Inserting
db.Exec("INSERT INTO users (name) VALUES ($1)", user.Name)
// Null/Unset → inserts NULL
// Valid → inserts the value
```

### DynamoDB Integration

Use the `nullddb` subpackage for DynamoDB:

```go
import "github.com/bjaus/null/nullddb"

type Item struct {
    PK   string                `dynamodbav:"pk"`
    Name nullddb.Value[string] `dynamodbav:"name"`
}

item := Item{
    PK:   "user#123",
    Name: nullddb.New("Alice"),
}

// Works with attributevalue.MarshalMap
av, _ := attributevalue.MarshalMap(item)
```

Convert between `null.Value` and `nullddb.Value`:

```go
// null.Value → nullddb.Value
nv := null.New("Alice")
dv := nullddb.From(nv)

// nullddb.Value embeds null.Value, so all methods work
dv.IsValid()  // true
dv.Get()      // "Alice"
```

## API Reference

### Constructors

| Function | Description |
|----------|-------------|
| `New[T](v)` | Create a valid Value |
| `NewNull[T]()` | Create a null Value |
| `NewPtr[T](p)` | Create from pointer (nil → null) |

### State Methods

| Method | Description |
|--------|-------------|
| `IsSet()` | True if null or valid |
| `IsNull()` | True if explicitly null |
| `IsValid()` | True if has a value |
| `State()` | Returns `Unset`, `Null`, or `Valid` |

### Value Methods

| Method | Description |
|--------|-------------|
| `Get()` | Returns value or zero |
| `GetOr(def)` | Returns value or default |
| `Ptr()` | Returns pointer or nil |

### Supported SQL Types

`string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`, `time.Time`, `[]byte`

## Design Decisions

**Why not pointers?**
- Can't distinguish unset from null
- Require heap allocation
- Awkward for literals (`&"hello"` doesn't work)

**Why not sql.NullString?**
- No unset/null distinction
- Bad JSON marshaling (`{"String": "x", "Valid": true}`)
- No generics

**Why a state enum instead of two bools?**
- Impossible to have invalid combinations
- Same memory after padding
- Clearer semantics

## License

MIT License - see [LICENSE](LICENSE) for details.
