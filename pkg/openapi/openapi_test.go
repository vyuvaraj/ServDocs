// Package openapi_test provides tests for the openapi package.
package openapi_test

import (
	"encoding/json"
	"os"
	"testing"

	"servdocs/pkg/openapi"
	"servdocs/pkg/parser"
)

func TestGenerate(t *testing.T) {
	doc := &parser.SrvDoc{
		Structs: []parser.StructDef{
			{Name: "Item", Description: "Simple item", Fields: []parser.StructField{{Name: "id", Type: "int"}}},
		},
		Routes: []parser.RouteDef{
			{Method: "POST", Path: "/item", InputType: "Item", OutputType: "string"},
		},
	}

	tmpJSON, err := os.CreateTemp("", "test_docs_*.json")
	if err != nil {
		t.Fatalf("failed to create temp json file: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	if err := openapi.Generate(doc, "Test API", tmpJSON.Name()); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	contentBytes, err := os.ReadFile(tmpJSON.Name())
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

func TestGenerateWithPathParams(t *testing.T) {
	doc := &parser.SrvDoc{
		Structs: []parser.StructDef{},
		Routes: []parser.RouteDef{
			{Method: "GET", Path: "/users/:id"},
			{Method: "DELETE", Path: "/users/:id/posts/:postId"},
		},
	}

	tmpJSON, err := os.CreateTemp("", "test_path_params_*.json")
	if err != nil {
		t.Fatalf("failed to create temp json file: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	if err := openapi.Generate(doc, "Path Params API", tmpJSON.Name()); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpJSON.Name())
	var parsed map[string]interface{}
	_ = json.Unmarshal(contentBytes, &parsed)

	paths := parsed["paths"].(map[string]interface{})
	if paths["/users/{id}"] == nil {
		t.Errorf("expected path /users/{id} to exist, got paths: %v", paths)
	}
	if paths["/users/{id}/posts/{postId}"] == nil {
		t.Errorf("expected path /users/{id}/posts/{postId} to exist, got paths: %v", paths)
	}
}
