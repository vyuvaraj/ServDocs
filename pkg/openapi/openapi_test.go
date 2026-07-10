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

func TestGenerateEmpty(t *testing.T) {
	doc := &parser.SrvDoc{}
	tmpJSON, err := os.CreateTemp("", "test_empty_*.json")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	if err := openapi.Generate(doc, "Empty API", tmpJSON.Name()); err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpJSON.Name())
	var parsed map[string]interface{}
	_ = json.Unmarshal(contentBytes, &parsed)

	paths := parsed["paths"].(map[string]interface{})
	if len(paths) != 0 {
		t.Errorf("expected empty paths, got: %v", paths)
	}
}

func TestGenerateMultipleRoutesSamePath(t *testing.T) {
	doc := &parser.SrvDoc{
		Routes: []parser.RouteDef{
			{Method: "GET", Path: "/user"},
			{Method: "POST", Path: "/user"},
		},
	}
	tmpJSON, err := os.CreateTemp("", "test_multiroute_*.json")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	if err := openapi.Generate(doc, "Multi API", tmpJSON.Name()); err != nil {
		t.Fatalf("failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpJSON.Name())
	var parsed map[string]interface{}
	_ = json.Unmarshal(contentBytes, &parsed)

	paths := parsed["paths"].(map[string]interface{})
	userPath := paths["/user"].(map[string]interface{})
	if userPath["get"] == nil || userPath["post"] == nil {
		t.Errorf("expected both get and post defined under /user path, got: %+v", userPath)
	}
}

func TestGenerateRouteWithNoInput(t *testing.T) {
	doc := &parser.SrvDoc{
		Routes: []parser.RouteDef{
			{Method: "POST", Path: "/ping"},
		},
	}
	tmpJSON, err := os.CreateTemp("", "test_noinput_*.json")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	if err := openapi.Generate(doc, "No Input API", tmpJSON.Name()); err != nil {
		t.Fatalf("failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpJSON.Name())
	var parsed map[string]interface{}
	_ = json.Unmarshal(contentBytes, &parsed)

	paths := parsed["paths"].(map[string]interface{})
	ping := paths["/ping"].(map[string]interface{})
	post := ping["post"].(map[string]interface{})
	if post["requestBody"] != nil {
		t.Errorf("expected no request body, got: %+v", post["requestBody"])
	}
}

func TestGenerateRouteWithNoOutput(t *testing.T) {
	doc := &parser.SrvDoc{
		Routes: []parser.RouteDef{
			{Method: "GET", Path: "/ping"},
		},
	}
	tmpJSON, err := os.CreateTemp("", "test_nooutput_*.json")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	if err := openapi.Generate(doc, "No Output API", tmpJSON.Name()); err != nil {
		t.Fatalf("failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpJSON.Name())
	var parsed map[string]interface{}
	_ = json.Unmarshal(contentBytes, &parsed)

	paths := parsed["paths"].(map[string]interface{})
	ping := paths["/ping"].(map[string]interface{})
	get := ping["get"].(map[string]interface{})
	responses := get["responses"].(map[string]interface{})
	resp200 := responses["200"].(map[string]interface{})
	if resp200["content"] != nil {
		t.Errorf("expected 200 response to have no content schema, got: %+v", resp200["content"])
	}
}

func TestGenerateTypeMappingBoolean(t *testing.T) {
	doc := &parser.SrvDoc{
		Structs: []parser.StructDef{
			{Name: "Data", Fields: []parser.StructField{{Name: "active", Type: "bool"}}},
		},
	}
	tmpJSON, err := os.CreateTemp("", "test_bool_*.json")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	if err := openapi.Generate(doc, "Bool API", tmpJSON.Name()); err != nil {
		t.Fatalf("failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpJSON.Name())
	var parsed map[string]interface{}
	_ = json.Unmarshal(contentBytes, &parsed)

	components := parsed["components"].(map[string]interface{})
	schemas := components["schemas"].(map[string]interface{})
	data := schemas["Data"].(map[string]interface{})
	props := data["properties"].(map[string]interface{})
	active := props["active"].(map[string]interface{})
	if active["type"] != "boolean" {
		t.Errorf("expected active property to map to 'boolean', got %q", active["type"])
	}
}

func TestGenerateTypeMappingInteger(t *testing.T) {
	doc := &parser.SrvDoc{
		Structs: []parser.StructDef{
			{Name: "Data", Fields: []parser.StructField{{Name: "count", Type: "int"}}},
		},
	}
	tmpJSON, err := os.CreateTemp("", "test_int_*.json")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	defer os.Remove(tmpJSON.Name())
	tmpJSON.Close()

	if err := openapi.Generate(doc, "Int API", tmpJSON.Name()); err != nil {
		t.Fatalf("failed: %v", err)
	}

	contentBytes, _ := os.ReadFile(tmpJSON.Name())
	var parsed map[string]interface{}
	_ = json.Unmarshal(contentBytes, &parsed)

	components := parsed["components"].(map[string]interface{})
	schemas := components["schemas"].(map[string]interface{})
	data := schemas["Data"].(map[string]interface{})
	props := data["properties"].(map[string]interface{})
	count := props["count"].(map[string]interface{})
	if count["type"] != "integer" {
		t.Errorf("expected count property to map to 'integer', got %q", count["type"])
	}
}
