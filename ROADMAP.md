# ServDocs Roadmap

This roadmap outlines the planned development phases for the ServDocs documentation generator.

---

## Differentiating Factors (Why ServDocs?)
* **Compiler-Aware Documentation**: Parses `.srv` source directly — route types, request/response schemas, middleware chains, and infrastructure declarations are all understood natively.
* **Dual Output**: Generates both interactive HTML documentation and OpenAPI 3.0 specs from a single parse pass.
* **Zero Configuration**: Point at a `.srv` file or directory and get publishable docs. No annotations or special comments required.

---

## Phase 1: Core Generator (Completed)
- [x] **`.srv` file parser** — Extracts routes, types, infrastructure declarations
- [x] **HTML documentation output** — Generates styled static HTML with route listing
- [x] **OpenAPI 3.0 output** — Generates JSON OpenAPI spec from parsed routes
- [x] **Local serve mode** — Built-in HTTP server for doc preview during development

## Phase 2: Quality & Testing (Pending — July 2026)

| # | Item | Effort | Description | Status |
|---|------|--------|-------------|--------|
| 2.1 | **Test suite** | Medium | Table-driven tests for parser, generator, and OpenAPI output. Currently zero tests exist | [ ] |
| 2.2 | **Dockerfile** | Small | Multi-stage build for containerized doc generation in CI pipelines | [ ] |
| 2.3 | **GitHub Actions CI** | Small | Automated build, test, and format checks | [ ] |
| 2.4 | **Multi-file project support** | Medium | Parse entire `serv.toml` projects, not just individual files. Cross-file type resolution | [ ] |
| 2.5 | **Markdown output format** | Small | Generate markdown docs alongside HTML — useful for GitHub wikis and README embedding | [ ] |

## Phase 3: Advanced Documentation (Pending)
- [ ] **Type schema rendering** — Render struct/interface definitions as expandable schema tables in HTML
- [ ] **Middleware chain documentation** — Show which middleware applies to which routes with order
- [ ] **Code examples in docs** — Include `.srv` usage examples alongside route documentation
- [ ] **Versioned docs** — Generate docs per git tag; host multiple versions side-by-side
- [ ] **Search** — Client-side full-text search across generated documentation
- [ ] **`serv docs serve --watch`** — File watcher that regenerates docs on `.srv` file changes

## Phase 4: Ecosystem Integration (Pending)
- [ ] **ServGate auto-registration** — Push generated OpenAPI specs to ServGate's auto-discovery endpoint
- [ ] **ServRegistry integration** — Include generated docs as package metadata on publish
- [ ] **ServConsole embedding** — Serve generated docs within the ServConsole documentation tab

> See [UNIFIED_ROADMAP.md](../servverse-repo/UNIFIED_ROADMAP.md) for the full ecosystem priority matrix.
