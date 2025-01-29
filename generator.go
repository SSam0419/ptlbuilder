package ptlbuilder

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"
)

// Generator handles protocol code generation
type Generator struct {
	spec *ProtocolSpec
}

// NewGenerator creates a new protocol generator
func NewGenerator(spec *ProtocolSpec) *Generator {
	return &Generator{
		spec: spec,
	}
}

// Generate generates protocol code and documentation
func (g *Generator) Generate() error {
	if err := g.spec.Validate(); err != nil {
		return fmt.Errorf("invalid protocol specification: %w", err)
	}

	// Generate protocol code
	protocolCode, err := g.generateProtocolCode()
	if err != nil {
		return fmt.Errorf("failed to generate protocol code: %w", err)
	}

	// Generate documentation
	docs, err := g.generateDocs()
	if err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	// Write files
	if err := g.writeFiles(protocolCode, docs); err != nil {
		return fmt.Errorf("failed to write files: %w", err)
	}

	return nil
}

func (g *Generator) generateProtocolCode() ([]byte, error) {
	funcMap := template.FuncMap{
		"lower": toLowerCamelCase,
	}

	tmpl, err := template.New("protocol").
		Funcs(funcMap).
		Parse(protocolTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse protocol template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.spec); err != nil {
		return nil, fmt.Errorf("failed to execute protocol template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to format generated code: %w", err)
	}

	return formatted, nil
}

func (g *Generator) generateDocs() ([]byte, error) {
	tmpl, err := template.New("doc").
		Funcs(template.FuncMap{
			"inc": func(i int) int {
				return i + 1
			},
			"char": func() string { return "```" },
		}).Parse(documentTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse documentation template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, g.spec); err != nil {
		return nil, fmt.Errorf("failed to execute documentation template: %w", err)
	}

	return buf.Bytes(), nil
}

func (g *Generator) writeFiles(protocolCode, docs []byte) error {

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	outputPath := filepath.Join(currentDir, g.spec.Package)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write protocol code
	protocolFile := filepath.Join(outputPath, fmt.Sprintf("%s.go", g.spec.Package))
	if err := os.WriteFile(protocolFile, protocolCode, 0644); err != nil {
		return fmt.Errorf("failed to write protocol file: %w", err)
	}

	// Write documentation
	docFile := filepath.Join(outputPath, "README.md")
	if err := os.WriteFile(docFile, docs, 0644); err != nil {
		return fmt.Errorf("failed to write documentation file: %w", err)
	}

	return nil
}

// toLowerCamelCase converts a PascalCase string to camelCase
func toLowerCamelCase(s string) string {
	if len(s) <= 1 {
		return s
	}
	return fmt.Sprintf("%c%s", toLower(s[0]), s[1:])
}

// toLower converts a byte to lowercase if it's an uppercase letter
func toLower(b byte) byte {
	return b | 0x20
}
