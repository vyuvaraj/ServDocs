package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
