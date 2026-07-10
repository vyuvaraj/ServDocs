// Package main integration tests — verify that the package-level wiring works
// end-to-end from the main package perspective. For detailed unit tests see
// pkg/parser, pkg/generator, and pkg/openapi.
package main

import (
	"os"
	"path/filepath"
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

func TestIntegration_CommandClientTS(t *testing.T) {
	_, srvPath := makeSampleDoc(t)
	doc, _ := parser.ParseSrvFile(srvPath)

	tmpDir, err := os.MkdirTemp("", "integration_sdk_ts")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := generator.GenerateClientSDK(doc, "typescript", tmpDir); err != nil {
		t.Fatalf("failed to generate ts sdk: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "client.ts"))
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if !strings.Contains(string(content), "interface User") {
		t.Errorf("expected TS SDK to contain interface User")
	}
}

func TestIntegration_CommandClientDart(t *testing.T) {
	_, srvPath := makeSampleDoc(t)
	doc, _ := parser.ParseSrvFile(srvPath)

	tmpDir, err := os.MkdirTemp("", "integration_sdk_dart")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := generator.GenerateClientSDK(doc, "dart", tmpDir); err != nil {
		t.Fatalf("failed to generate dart sdk: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "client.dart"))
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if !strings.Contains(string(content), "class User") {
		t.Errorf("expected Dart SDK to contain class User")
	}
}

func TestIntegration_CommandClientSwift(t *testing.T) {
	_, srvPath := makeSampleDoc(t)
	doc, _ := parser.ParseSrvFile(srvPath)

	tmpDir, err := os.MkdirTemp("", "integration_sdk_swift")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := generator.GenerateClientSDK(doc, "swift", tmpDir); err != nil {
		t.Fatalf("failed to generate swift sdk: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "client.swift"))
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if !strings.Contains(string(content), "struct User: Codable") {
		t.Errorf("expected Swift SDK to contain struct User")
	}
}

func TestIntegration_CommandOpenAPI(t *testing.T) {
	_, srvPath := makeSampleDoc(t)
	doc, _ := parser.ParseSrvFile(srvPath)

	tmpFile, err := os.CreateTemp("", "integration_openapi_*.json")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	if err := openapi.Generate(doc, "E2E API", tmpFile.Name()); err != nil {
		t.Fatalf("failed: %v", err)
	}

	content, _ := os.ReadFile(tmpFile.Name())
	if !strings.Contains(string(content), `"openapi": "3.0.3"`) {
		t.Errorf("expected valid OpenAPI version header")
	}
}
