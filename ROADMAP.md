# ServDocs Roadmap

This roadmap outlines the planned development phases for the ServDocs documentation generator.

---

## Differentiating Factors (Why ServDocs?)
* **Compiler-Aware Documentation**: Parses `.srv` source directly — route types, request/response schemas, middleware chains, and infrastructure declarations are all understood natively.
* **Dual Output**: Generates both interactive HTML documentation and OpenAPI 3.0 specs from a single parse pass.
* **Multi-Language SDKs**: Auto-generates typed client libraries for TypeScript, Dart, and Swift.
* **Zero Configuration**: Point at a `.srv` file and get publishable docs. No annotations or special comments required.

---

## Phase 1: Core Generator (Completed)
- [x] **`.srv` file parser** — Extracts routes, types, and function declarations
- [x] **HTML documentation output** — Generates styled static HTML with route listing
- [x] **OpenAPI 3.0 output** — Generates JSON OpenAPI spec from parsed routes
- [x] **Local serve mode** — Built-in HTTP server for doc preview during development

## Phase 2: Quality & Testing (Completed — July 2026)

| # | Item | Effort | Description | Status |
|---|------|--------|-------------|--------|
| 2.1 | **Test suite** | Medium | Table-driven tests for parser, generator, and OpenAPI output | [x] |
| 2.2 | **Dockerfile** | Small | Multi-stage build for containerized doc generation in CI pipelines | [x] |
| 2.3 | **GitHub Actions CI** | Small | Automated build, test, and format checks | [x] |
| 2.4 | **Multi-file project support** | Medium | Parse entire `serv.toml` projects, not just individual files. Cross-file type resolution | [ ] |
| 2.5 | **Markdown output format** | Small | Generate markdown docs alongside HTML — useful for GitHub wikis and README embedding | [ ] |

## Phase 3: Advanced Documentation (Completed)
- [x] **Type schema rendering** — Render struct/interface definitions as expandable schema tables in HTML
- [x] **Middleware chain documentation** — Show which middleware applies to which routes with order [July 9, 2026]
- [x] **Versioned docs** — Generate docs per git tag; host multiple versions side-by-side with version selector
- [x] **Search** — Client-side full-text search with highlight across generated documentation
- [x] **Multi-language code examples** — cURL, Go, and JavaScript snippets auto-generated per route
- [ ] **Code examples in docs** — Include `.srv` usage examples alongside route documentation

## Phase 4: Ecosystem Integration (Pending)
- [ ] **ServGate auto-registration** — Push generated OpenAPI specs to ServGate's auto-discovery endpoint
- [ ] **ServRegistry integration** — Include generated docs as package metadata on publish
- [ ] **ServConsole embedding** — Serve generated docs within the ServConsole documentation tab

> See [UNIFIED_ROADMAP.md](../servverse-repo/UNIFIED_ROADMAP.md) for the full ecosystem priority matrix.


---

## Phase 5: Code Health & Test Coverage (Completed — July 2026, Phase 22 QC.5)

| # | Item | Effort | Description | Status |
|---|------|--------|-------------|--------|
| 5.1 | **Add pkg/ structure** | Medium | Created `pkg/parser/`, `pkg/generator/`, `pkg/openapi/` with clean interfaces. `main.go` reduced to <100 lines | [x] |
| 5.2 | **Expand test suite** | Medium | Grew from 5 → 17 test functions across 4 packages: parser accuracy, OpenAPI validation, path param normalization, HTML generation, versioned docs, all three SDK targets, unsupported lang error path, plus 2 integration tests | [x] |
