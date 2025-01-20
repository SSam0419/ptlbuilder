package ptlbuilder

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

func Generate(p ProtocolSpec) {
	log.Println("start generating ..")
	spec := p

	funcMap := template.FuncMap{
		"lower": func(s string) string {
			// Convert field name to parameter name (e.g., ClientAddr -> clientAddr)
			if len(s) <= 1 {
				return s
			}
			return fmt.Sprintf("%s%s", string(s[0]|32), s[1:])
		},
	}

	tmpl := template.Must(template.New("protocol").Funcs(funcMap).Parse(protocolTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, spec); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		return
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		fmt.Printf("Error formatting code: %v\n", err)
		return
	}

	// Write to file
	currentDir, err := os.Getwd()
	if err != nil {
		log.Printf("Error getting current directory: %v\n", err)
		return
	}

	outputPath := filepath.Join(currentDir, "protocol")

	// Ensure the directory exists
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		log.Printf("Error creating directory: %v\n", err)
		return
	}

	// Write the file
	filePath := filepath.Join(outputPath, "protocol.go")
	if err := os.WriteFile(filePath, formatted, 0644); err != nil {
		log.Printf("Error writing file: %v\n", err)
		return
	}

	log.Printf("Successfully generated: %s\n", filePath)
}
