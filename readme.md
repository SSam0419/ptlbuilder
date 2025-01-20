# PTLBuilder

PTLBuilder is a Go library for generating protocol code. It simplifies the creation of communication protocols by automatically generating the necessary Go code based on a protocol specification.

## Installation

```bash
go get github.com/SSam0419/ptlbuilder
```

## Usage

```go
package main

import "github.com/SSam0419/ptlbuilder"

func main() {
    ptlbuilder.Generate(ptlbuilder.ProtocolSpec{
        Package: "protocol",
        Commands: []ptlbuilder.Command{
            {
                Name: "RegisterClient",
                Fields: []ptlbuilder.Field{
                    {Name: "ClientAddr", Type: "string"},
                },
            },
            {
                Name: "ListenTopic",
                Fields: []ptlbuilder.Field{
                    {Name: "Topic", Type: "string"},
                    {Name: "ClientAddr", Type: "string"},
                },
            },
            {
                Name: "SendMessage",
                Fields: []ptlbuilder.Field{
                    {Name: "Topic", Type: "string"},
                    {Name: "ClientAddr", Type: "string"},
                    {Name: "Payload", Type: "[]byte"},
                },
            },
        },
    })
}
```


## Protocol Specification

### ProtocolSpec
The main structure for defining your protocol:

```go
type ProtocolSpec struct {
    Package  string    // Target package name
    Commands []Command // List of commands
}
```

### Command
Defines a single command in the protocol:

```go
type Command struct {
    Name   string  // Command name
    Fields []Field // Command fields
}
```

### Field
Defines a field within a command:

```go
type Field struct {
    Name string // Field name
    Type string // Field type (Go type)
}
```

### Generated Code
The library will generate a protocol.go file in the protocol directory