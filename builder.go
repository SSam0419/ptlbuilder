// Package ptlbuilder provides utilities for building protocol specifications.
package ptlbuilder

import (
	"errors"
	"strings"
)

// Common errors
var (
	ErrEmptyName    = errors.New("name cannot be empty")
	ErrEmptyType    = errors.New("type cannot be empty")
	ErrEmptyPackage = errors.New("package name cannot be empty")
	ErrZeroTimeout  = errors.New("timeout must be greater than 0")
)

// Field represents a single field in a protocol command.
type Field struct {
	Name string
	Type string
}

// NewField creates a new field with validation.
func NewField(name, fieldType string) *Field {
	if strings.TrimSpace(name) == "" {
		panic(ErrEmptyName)
	}
	if strings.TrimSpace(fieldType) == "" {
		panic(ErrEmptyType)
	}

	return &Field{
		Name: name,
		Type: fieldType,
	}
}

// Command represents a protocol command with its fields.
type Command struct {
	Name   string
	Fields []Field
}

// NewCommand creates a new command with validation.
// Only allow UpperCamel / Pascal Case
func NewCommand(name string) *Command {
	if strings.TrimSpace(name) == "" {
		panic(ErrEmptyName)
	}
	if !isPascalCase(name) {
		panic("command name must be in PascalCase")
	}

	return &Command{
		Name:   name,
		Fields: make([]Field, 0),
	}
}

// isPascalCase checks if a string is in PascalCase.
func isPascalCase(s string) bool {
	if s == "" {
		return false
	}
	if !strings.HasPrefix(s, strings.ToUpper(string(s[0]))) {
		return false
	}
	for _, r := range s[1:] {
		if r == '_' || r == '-' || r == ' ' {
			return false
		}
	}
	return true
}

// AddField adds a new field to the command with validation.
func (c *Command) AddField(name string, fieldType string) *Command {
	field := NewField(name, fieldType)
	c.Fields = append(c.Fields, *field)
	return c
}

// ProtocolSpec represents the complete protocol specification.
type ProtocolSpec struct {
	Package  string
	Commands []*Command
	Timeout  uint // maximum seconds allowed to decode message
}

// NewProtocolSpec creates a new protocol specification with validation.
func NewProtocolSpec(pkg string, timeout uint) *ProtocolSpec {
	if strings.TrimSpace(pkg) == "" {
		panic(ErrEmptyPackage)
	}
	if timeout == 0 {
		panic(ErrZeroTimeout)
	}

	return &ProtocolSpec{
		Package:  pkg,
		Commands: make([]*Command, 0),
		Timeout:  timeout,
	}
}

// AddCommand adds a new command to the protocol specification.
func (p *ProtocolSpec) AddCommand(cmd *Command) *ProtocolSpec {
	p.Commands = append(p.Commands, cmd)
	return p
}

// Validate checks if the protocol specification is valid.
func (p *ProtocolSpec) Validate() error {
	if strings.TrimSpace(p.Package) == "" {
		return ErrEmptyPackage
	}
	if p.Timeout == 0 {
		return ErrZeroTimeout
	}

	for _, cmd := range p.Commands {
		if strings.TrimSpace(cmd.Name) == "" {
			return ErrEmptyName
		}
		for _, field := range cmd.Fields {
			if strings.TrimSpace(field.Name) == "" {
				return ErrEmptyName
			}
			if strings.TrimSpace(field.Type) == "" {
				return ErrEmptyType
			}
		}
	}

	return nil
}
