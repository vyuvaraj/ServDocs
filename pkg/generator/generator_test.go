// Package generator_test provides tests for the generator package.
package generator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"servdocs/pkg/generator"
	"servdocs/pkg/parser"
)

func TestGenerateHtml(t *testing.T) {
	doc := &parser.SrvDoc{
		Structs: []parser.StructDef{
			{Name: "Item", Description: "Simple item", Fields: []parser.StructField{{Name: "id", Type: "int"}}},
		},
		Routes: []parser.RouteDef{
			{Method: "POST", Path: "/item", InputType: "Item", OutputType: "string"},
		},
	}

	tmpHtml, err := os.CreateTemp("", "test_docs_*.html")
	if err != nil {
		t.Fatalf("failed to create temp html file: %v", err)
	}
	defer os.Remove(tmpHtml.Name())
	tmpHtml.Close()

	if err := generator.GenerateHtml(doc, "Test Title", tmpHtml.Name(), "", ""); err != nil {
		t.Fatalf("GenerateHtml failed: %v", err)
	}

	contentBytes, err := os.ReadFile(tmpHtml.Name())
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

func TestGenerateHtmlVersioned(t *testing.T) {
	doc := &parser.SrvDoc{
		Structs: []parser.StructDef{
			{Name: "Item", Description: "Simple item", Fields: []parser.StructField{{Name: "id", Type: "int"}, {Name: "name", Type: "string"}}},
		},
		Routes: []parser.RouteDef{
			{Method: "POST", Path: "/item", InputType: "Item", OutputType: "Item"},
		},
	}

	tmpDir, err := os.MkdirTemp("", "test_docs_versioned")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	v1Dir := filepath.Join(tmpDir, "v1.0.0")
	v2Dir := filepath.Join(tmpDir, "v2.0.0")
	if err := os.MkdirAll(v1Dir, 0755); err != nil {
		t.Fatalf("failed to create v1 dir: %v", err)
	}
	if err := os.MkdirAll(v2Dir, 0755); err != nil {
		t.Fatalf("failed to create v2 dir: %v", err)
	}

	if err := generator.GenerateHtml(doc, "Test Title", filepath.Join(v1Dir, "index.html"), tmpDir, "v1.0.0"); err != nil {
		t.Fatalf("GenerateHtml failed for v1.0.0: %v", err)
	}
	if err := generator.GenerateHtml(doc, "Test Title", filepath.Join(v2Dir, "index.html"), tmpDir, "v2.0.0"); err != nil {
		t.Fatalf("GenerateHtml failed for v2.0.0: %v", err)
	}

	content1, err := os.ReadFile(filepath.Join(v1Dir, "index.html"))
	if err != nil {
		t.Fatalf("failed to read v1 index: %v", err)
	}
	content2, err := os.ReadFile(filepath.Join(v2Dir, "index.html"))
	if err != nil {
		t.Fatalf("failed to read v2 index: %v", err)
	}

	if !strings.Contains(string(content1), "v1.0.0") || !strings.Contains(string(content1), "v2.0.0") {
		t.Errorf("v1 index doesn't list all versions: %s", string(content1))
	}
	if !strings.Contains(string(content2), "v1.0.0") || !strings.Contains(string(content2), "v2.0.0") {
		t.Errorf("v2 index doesn't list all versions: %s", string(content2))
	}
}

func TestGenerateClientSDK(t *testing.T) {
	doc := &parser.SrvDoc{
		Structs: []parser.StructDef{
			{Name: "User", Description: "Simple user", Fields: []parser.StructField{{Name: "id", Type: "int"}, {Name: "username", Type: "string"}}},
		},
		Routes: []parser.RouteDef{
			{Method: "POST", Path: "/user", InputType: "User", OutputType: "User"},
		},
	}

	tmpDir, err := os.MkdirTemp("", "test_sdk")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. Test TypeScript
	if err := generator.GenerateClientSDK(doc, "typescript", tmpDir); err != nil {
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
	if err := generator.GenerateClientSDK(doc, "dart", tmpDir); err != nil {
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
	if err := generator.GenerateClientSDK(doc, "swift", tmpDir); err != nil {
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

func TestGenerateClientSDKUnsupported(t *testing.T) {
	doc := &parser.SrvDoc{}
	tmpDir, err := os.MkdirTemp("", "test_sdk_bad")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err := generator.GenerateClientSDK(doc, "ruby", tmpDir); err == nil {
		t.Errorf("expected error for unsupported language 'ruby', got nil")
	}
}
