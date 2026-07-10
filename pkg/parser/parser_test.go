// Package parser_test provides tests for the parser package.
package parser_test

import (
	"os"
	"testing"

	"servdocs/pkg/parser"
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

	doc, err := parser.ParseSrvFile(tmpfile.Name())
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
