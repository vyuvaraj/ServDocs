// Package parser_test provides tests for the parser package.
package parser_test

import (
	"os"
	"testing"

	"servdocs/pkg/parser"
)

func createTempSrvFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "test_docs_*.srv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	t.Cleanup(func() { os.Remove(tmpfile.Name()) })

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpfile.Close()
	return tmpfile.Name()
}

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
	path := createTempSrvFile(t, srvContent)
	doc, err := parser.ParseSrvFile(path)
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

func TestParseEmptyFile(t *testing.T) {
	path := createTempSrvFile(t, "")
	doc, err := parser.ParseSrvFile(path)
	if err != nil {
		t.Fatalf("failed to parse empty file: %v", err)
	}
	if len(doc.Structs) != 0 || len(doc.Routes) != 0 || len(doc.Functions) != 0 {
		t.Errorf("expected empty doc, got %+v", doc)
	}
}

func TestParseOnlyComments(t *testing.T) {
	path := createTempSrvFile(t, "/// some comment\n/// another comment\n")
	doc, err := parser.ParseSrvFile(path)
	if err != nil {
		t.Fatalf("failed to parse comments: %v", err)
	}
	if len(doc.Structs) != 0 || len(doc.Routes) != 0 || len(doc.Functions) != 0 {
		t.Errorf("expected empty doc, got %+v", doc)
	}
}

func TestParseMultipleStructs(t *testing.T) {
	srvContent := `
struct A {
    a: int
}
struct B {
    b: string
}
`
	path := createTempSrvFile(t, srvContent)
	doc, err := parser.ParseSrvFile(path)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	if len(doc.Structs) != 2 {
		t.Errorf("expected 2 structs, got %d", len(doc.Structs))
	}
}

func TestParseStructWithComments(t *testing.T) {
	srvContent := `
/// Struct comment
struct Test {
    /// Field comment
    field: string
}
`
	path := createTempSrvFile(t, srvContent)
	doc, err := parser.ParseSrvFile(path)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(doc.Structs) != 1 || doc.Structs[0].Description != "Struct comment" {
		t.Errorf("missing struct comment, got: %q", doc.Structs[0].Description)
	}
	if doc.Structs[0].Fields[0].Description != "Field comment" {
		t.Errorf("missing field comment, got: %q", doc.Structs[0].Fields[0].Description)
	}
}

func TestParseRouteWithNoInputOrOutput(t *testing.T) {
	srvContent := `
route GET /ping ()
`
	path := createTempSrvFile(t, srvContent)
	doc, err := parser.ParseSrvFile(path)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(doc.Routes) != 1 {
		t.Fatalf("expected 1 route, got 0")
	}
	route := doc.Routes[0]
	if route.InputType != "" || route.OutputType != "" {
		t.Errorf("expected empty input/output, got input=%q output=%q", route.InputType, route.OutputType)
	}
}

func TestParseRouteWithMiddlewares(t *testing.T) {
	srvContent := `
route POST /submit (Data) use [validate, rate_limit, cors] -> Response
`
	path := createTempSrvFile(t, srvContent)
	doc, err := parser.ParseSrvFile(path)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(doc.Routes) != 1 {
		t.Fatalf("expected 1 route")
	}
	route := doc.Routes[0]
	expected := []string{"validate", "rate_limit", "cors"}
	if len(route.Middlewares) != len(expected) {
		t.Fatalf("expected %d middlewares, got %d", len(expected), len(route.Middlewares))
	}
	for i, m := range expected {
		if route.Middlewares[i] != m {
			t.Errorf("expected middleware %q, got %q", m, route.Middlewares[i])
		}
	}
}

func TestParseFunctionsNoReturn(t *testing.T) {
	srvContent := `
fn doSomething(x, y)
`
	path := createTempSrvFile(t, srvContent)
	doc, err := parser.ParseSrvFile(path)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	if len(doc.Functions) != 1 {
		t.Fatalf("expected 1 function")
	}
	fn := doc.Functions[0]
	if fn.ReturnType != "" {
		t.Errorf("expected empty return type, got %q", fn.ReturnType)
	}
}

func TestParseInvalidSyntaxNoPanic(t *testing.T) {
	srvContent := `
invalid syntax structure here {()
`
	path := createTempSrvFile(t, srvContent)
	_, err := parser.ParseSrvFile(path)
	if err != nil {
		t.Fatalf("expected parsing to complete without error, got: %v", err)
	}
}
