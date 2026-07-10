// Package main integration tests — verify that the package-level wiring works
// end-to-end from the main package perspective. For detailed unit tests see
// pkg/parser, pkg/generator, and pkg/openapi.
package main

import (
	"os"
	"strings"
	"testing"

	"servdocs/pkg/generator"
	"servdocs/pkg/openapi"
	"servdocs/pkg/parser"
)

func makeSampleDoc(t *testing.T) (*parser.SrvDoc, string) {
	t.Helper()
	srvContent := `
/// User information
struct User {
    /// The unique username
    username: string
    age: int
}

/// Retrieve user profile
route GET /user (id) -> User

/// Secure data endpoint
route "POST" "/secure-data" (req) use [auth, logging]

/// Process payment action
fn processPayment(amount, method) -> string
`
	tmpfile, err := os.CreateTemp("", "integration_*.srv")
	if err != nil {
		t.Fatalf("failed to create temp srv file: %v", err)
	}
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })

	if _, err := tmpfile.Write([]byte(srvContent)); err != nil {
		t.Fatalf("failed to write temp srv file: %v", err)
	}
	tmpfile.Close()
	return nil, tmpfile.Name()
}

func TestIntegration_ParseAndGenerateHtml(t *testing.T) {
	_, srvPath := makeSampleDoc(t)

	doc, err := parser.ParseSrvFile(srvPath)
	if err != nil {
		t.Fatalf("ParseSrvFile: %v", err)
	}
	if len(doc.Structs) == 0 || len(doc.Routes) == 0 {
		t.Fatalf("expected non-empty doc, got %+v", doc)
	}

	tmpHTML, err := os.CreateTemp("", "integration_*.html")
	if err != nil {
		t.Fatalf("failed to create temp html: %v", err)
	}
	t.Cleanup(func() { os.Remove(tmpHTML.Name()) })
	tmpHTML.Close()

	if err := generator.GenerateHtml(doc, "Integration Test", tmpHTML.Name(), "", ""); err != nil {
		t.Fatalf("GenerateHtml: %v", err)
	}

	content, _ := os.ReadFile(tmpHTML.Name())
	if !strings.Contains(string(content), "Integration Test") {
		t.Errorf("generated HTML missing title")
	}
}

func TestIntegration_ParseAndGenerateOpenAPI(t *testing.T) {
	_, srvPath := makeSampleDoc(t)

	doc, err := parser.ParseSrvFile(srvPath)
	if err != nil {
		t.Fatalf("ParseSrvFile: %v", err)
	}

	tmpJSON, err := os.CreateTemp("", "integration_*.json")
	if err != nil {
		t.Fatalf("failed to create temp json: %v", err)
	}
	t.Cleanup(func() { os.Remove(tmpJSON.Name()) })
	tmpJSON.Close()

	if err := openapi.Generate(doc, "Integration API", tmpJSON.Name()); err != nil {
		t.Fatalf("openapi.Generate: %v", err)
	}

	content, _ := os.ReadFile(tmpJSON.Name())
	if !strings.Contains(string(content), "Integration API") {
		t.Errorf("generated OpenAPI JSON missing title")
	}
}
