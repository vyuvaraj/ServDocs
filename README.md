# ServDocs — Auto-Generated Documentation

ServDocs is a CLI tool that parses Serv-lang `.srv` files to auto-generate static HTML documentation, OpenAPI 3.0 specifications, and multi-language client SDKs.

## Package Structure

```
ServDocs/
├── main.go                     # CLI entry point (<100 lines)
├── main_test.go                # Integration tests
└── pkg/
    ├── parser/                 # .srv file parsing & AST types
    │   ├── parser.go
    │   └── parser_test.go
    ├── generator/              # HTML documentation & client SDK generation
    │   ├── generator.go        # GenerateHtml — self-contained HTML site
    │   ├── client.go           # GenerateClientSDK — TypeScript, Dart, Swift
    │   └── generator_test.go
    └── openapi/                # OpenAPI 3.0 spec generation
        ├── openapi.go
        └── openapi_test.go
```

## CLI Usage

```bash
# Generate a static HTML documentation site
servdocs generate --input example.srv --output docs.html --title "My Service API"

# Generate with versioning (creates out-dir/v1.0.0/index.html)
servdocs generate --input example.srv --out-dir ./docs --version-tag v1.0.0

# Generate an OpenAPI 3.0 JSON specification
servdocs openapi --input example.srv --output openapi.json --title "My OpenAPI Spec"

# Serve docs locally (live preview)
servdocs serve --input example.srv --port 3000

# Generate a typed client SDK
servdocs client --input example.srv --lang typescript --output ./sdk
servdocs client --input example.srv --lang dart       --output ./sdk
servdocs client --input example.srv --lang swift      --output ./sdk
```

## Running Tests

```bash
go test ./...
```

All four packages are independently testable:

```
ok  servdocs                  (integration)
ok  servdocs/pkg/parser       (unit)
ok  servdocs/pkg/generator    (unit)
ok  servdocs/pkg/openapi      (unit)
```
