# ServDocs — Auto-Generated Documentation

ServDocs is a CLI tool that parses Serv-lang `.srv` files to auto-generate static HTML documentation and OpenAPI 3.0 specifications.

## CLI Usage

```bash
servdocs generate -input example.srv -output docs.html -title "My Service API"
servdocs openapi -input example.srv -output openapi.json -title "My OpenAPI Spec"
servdocs serve -input example.srv -port 3000
```
