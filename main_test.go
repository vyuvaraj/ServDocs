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

/// Process payment action
fn processPayment(amount, method) -> string
`
	tmpfile, err := ioutil.TempFile("", "test_docs_*.srv")
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

	if len(doc.Routes) != 1 || doc.Routes[0].Path != "/user" || doc.Routes[0].Method != "GET" {
		t.Errorf("expected route GET /user, got: %+v", doc.Routes)
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

	err = GenerateHtml(doc, "Test Title", tmpHtml.Name())
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

