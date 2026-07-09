package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSrvFile(t *testing.T) {
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
	tmpfile, err := os.CreateTemp("", "test_docs_*.srv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(srvContent)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpfile.Close()

	doc, err := ParseSrvFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseSrvFile failed: %v", err)
	}

	if len(doc.Structs) != 1 || doc.Structs[0].Name != "User" {
		t.Errorf("expected struct 'User', got: %+v", doc.Structs)
	}
	if len(doc.Structs[0].Fields) != 2 || doc.Structs[0].Fields[0].Name != "username" {
		t.Errorf("expected 2 struct fields, got: %+v", doc.Structs[0].Fields)
	}

	if len(doc.Routes) != 2 {
		t.Errorf("expected 2 routes, got: %d", len(doc.Routes))
	} else {
		if doc.Routes[0].Path != "/user" || doc.Routes[0].Method != "GET" {
			t.Errorf("expected first route GET /user, got: %+v", doc.Routes[0])
		}
		if doc.Routes[1].Path != "/secure-data" || doc.Routes[1].Method != "POST" {
			t.Errorf("expected second route POST /secure-data, got: %+v", doc.Routes[1])
		}
		expectedMiddlewares := []string{"auth", "logging"}
		if len(doc.Routes[1].Middlewares) != 2 || doc.Routes[1].Middlewares[0] != "auth" || doc.Routes[1].Middlewares[1] != "logging" {
			t.Errorf("expected middlewares %v, got: %v", expectedMiddlewares, doc.Routes[1].Middlewares)
		}
	}

	if len(doc.Functions) != 1 || doc.Functions[0].Name != "processPayment" {
		t.Errorf("expected function 'processPayment', got: %+v", doc.Functions)
	}
}

func TestGenerateHtml(t *testing.T) {
	doc := &SrvDoc{
		Structs: []StructDef{
			{Name: "Item", Description: "Simple item", Fields: []StructField{{Name: "id", Type: "int"}}},
		},
		Routes: []RouteDef{
			{Method: "POST", Path: "/item", InputType: "Item", OutputType: "string"},
		},
	}

	tmpHtml, err := ioutil.TempFile("", "test_docs_*.html")
	if err != nil {
		t.Fatalf("failed to create temp html file: %v", err)
	}
	defer os.Remove(tmpHtml.Name())
	tmpHtml.Close()

	err = GenerateHtml(doc, "Test Title", tmpHtml.Name(), "", "")
	if err != nil {
		t.Fatalf("GenerateHtml failed: %v", err)
	}

	contentBytes, err := ioutil.ReadFile(tmpHtml.Name())
	if err != nil {
		t.Fatalf("failed to read generated html: %v", err)
	}
	content := string(contentBytes)

	if !strings.Contains(content, "Test Title") {
		t.Errorf("HTML missing title 'Test Title'")
	}
	if !strings.Contains(content, "/item") {
		t.Errorf("HTML missing route path '/item'")
	}
	if !strings.Contains(content, "srvSchemas =") {
		t.Errorf("HTML missing serialized schemas JSON 'srvSchemas = ...'")
	}
	if !strings.Contains(content, "Item") {
		t.Errorf("HTML missing struct name 'Item' in serialized schema")
	}
	if !strings.Contains(content, "toggle-schema-btn") {
		t.Errorf("HTML missing toggle schema button element class 'toggle-schema-btn'")
	}
	if !strings.Contains(content, "search-status") {
		t.Errorf("HTML missing search results status element ID 'search-status'")
	}
}

func TestGenerateOpenAPI(t *testing.T) {
	doc := &SrvDoc{
		Structs: []StructDef{
			{Name: "Item", Description: "Simple item", Fields: []StructField{{Name: "id", Type: "int"}}},
		},
		Routes: []RouteDef{
			{Method: "POST", Path: "/item", InputType: "Item", OutputType: "string"},
		},
	}

	tmpJSON, err := ioutil.TempFile("", "test_docs_*.json")
	if err != nil {
		t.Fatalf("failed to create temp json file: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	err = GenerateOpenAPI(doc, "Test API", tmpJSON.Name())
	if err != nil {
		t.Fatalf("GenerateOpenAPI failed: %v", err)
	}

	contentBytes, err := ioutil.ReadFile(tmpJSON.Name())
	if err != nil {
		t.Fatalf("failed to read generated OpenAPI: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(contentBytes, &parsed); err != nil {
		t.Fatalf("failed to unmarshal OpenAPI JSON: %v", err)
	}

	info := parsed["info"].(map[string]interface{})
	if info["title"] != "Test API" {
		t.Errorf("expected OpenAPI title 'Test API', got %v", info["title"])
	}

	paths := parsed["paths"].(map[string]interface{})
	if paths["/item"] == nil {
		t.Errorf("missing path '/item' in OpenAPI specs")
	}
}

func TestGenerateClientSDK(t *testing.T) {
	doc := &SrvDoc{
		Structs: []StructDef{
			{Name: "User", Description: "Simple user", Fields: []StructField{{Name: "id", Type: "int"}, {Name: "username", Type: "string"}}},
		},
		Routes: []RouteDef{
			{Method: "POST", Path: "/user", InputType: "User", OutputType: "User"},
		},
	}

	tmpDir, err := os.MkdirTemp("", "test_sdk")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. Test TypeScript
	err = GenerateClientSDK(doc, "typescript", tmpDir)
	if err != nil {
		t.Fatalf("typescript generation failed: %v", err)
	}
	tsBytes, err := os.ReadFile(filepath.Join(tmpDir, "client.ts"))
	if err != nil {
		t.Fatalf("failed to read ts: %v", err)
	}
	tsCode := string(tsBytes)
	if !strings.Contains(tsCode, "export interface User") || !strings.Contains(tsCode, "class APIClient") {
		t.Errorf("invalid typescript output: %s", tsCode)
	}

	// 2. Test Dart
	err = GenerateClientSDK(doc, "dart", tmpDir)
	if err != nil {
		t.Fatalf("dart generation failed: %v", err)
	}
	dartBytes, err := os.ReadFile(filepath.Join(tmpDir, "client.dart"))
	if err != nil {
		t.Fatalf("failed to read dart: %v", err)
	}
	dartCode := string(dartBytes)
	if !strings.Contains(dartCode, "class User") || !strings.Contains(dartCode, "fromJson") {
		t.Errorf("invalid dart output: %s", dartCode)
	}

	// 3. Test Swift
	err = GenerateClientSDK(doc, "swift", tmpDir)
	if err != nil {
		t.Fatalf("swift generation failed: %v", err)
	}
	swiftBytes, err := os.ReadFile(filepath.Join(tmpDir, "client.swift"))
	if err != nil {
		t.Fatalf("failed to read swift: %v", err)
	}
	swiftCode := string(swiftBytes)
	if !strings.Contains(swiftCode, "public struct User: Codable") {
		t.Errorf("invalid swift output: %s", swiftCode)
	}
}

func TestGenerateHtmlVersioned(t *testing.T) {
	doc := &SrvDoc{
		Structs: []StructDef{
			{Name: "Item", Description: "Simple item", Fields: []StructField{{Name: "id", Type: "int"}, {Name: "name", Type: "string"}}},
		},
		Routes: []RouteDef{
			{Method: "POST", Path: "/item", InputType: "Item", OutputType: "Item"},
		},
	}

	tmpDir, err := os.MkdirTemp("", "test_docs_versioned")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create directories for two versions
	v1Dir := filepath.Join(tmpDir, "v1.0.0")
	v2Dir := filepath.Join(tmpDir, "v2.0.0")
	if err := os.MkdirAll(v1Dir, 0755); err != nil {
		t.Fatalf("failed to create v1 dir: %v", err)
	}
	if err := os.MkdirAll(v2Dir, 0755); err != nil {
		t.Fatalf("failed to create v2 dir: %v", err)
	}

	err = GenerateHtml(doc, "Test Title", filepath.Join(v1Dir, "index.html"), tmpDir, "v1.0.0")
	if err != nil {
		t.Fatalf("GenerateHtml failed for v1.0.0: %v", err)
	}

	err = GenerateHtml(doc, "Test Title", filepath.Join(v2Dir, "index.html"), tmpDir, "v2.0.0")
	if err != nil {
		t.Fatalf("GenerateHtml failed for v2.0.0: %v", err)
	}

	// Verify both exist
	content1, err := os.ReadFile(filepath.Join(v1Dir, "index.html"))
	if err != nil {
		t.Fatalf("failed to read v1 index: %v", err)
	}
	content2, err := os.ReadFile(filepath.Join(v2Dir, "index.html"))
	if err != nil {
		t.Fatalf("failed to read v2 index: %v", err)
	}

	// Check if version selector with both versions is present
	if !strings.Contains(string(content1), "v1.0.0") || !strings.Contains(string(content1), "v2.0.0") {
		t.Errorf("v1 index doesn't list all versions: %s", string(content1))
	}
	if !strings.Contains(string(content2), "v1.0.0") || !strings.Contains(string(content2), "v2.0.0") {
		t.Errorf("v2 index doesn't list all versions: %s", string(content2))
	}
}


