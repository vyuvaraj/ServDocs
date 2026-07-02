package main

import (
	"encoding/json"
	"os"
	"strings"
)

type OpenAPISpec struct {
	OpenAPI    string                 `json:"openapi"`
	Info       map[string]interface{} `json:"info"`
	Paths      map[string]interface{} `json:"paths"`
	Components map[string]interface{} `json:"components"`
}

func GenerateOpenAPI(doc *SrvDoc, title, outputPath string) error {
	spec := OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: map[string]interface{}{
			"title":       title,
			"version":     "1.0.0",
			"description": "Auto-generated OpenAPI specification from Servverse .srv source files",
		},
		Paths: make(map[string]interface{}),
		Components: map[string]interface{}{
			"schemas": make(map[string]interface{}),
		},
	}

	schemas := spec.Components["schemas"].(map[string]interface{})

	// Map Structs to Schemas
	for _, s := range doc.Structs {
		properties := make(map[string]interface{})
		for _, f := range s.Fields {
			pType := "string"
			fType := strings.ToLower(f.Type)
			switch fType {
			case "int", "integer":
				pType = "integer"
			case "bool", "boolean":
				pType = "boolean"
			}
			properties[f.Name] = map[string]interface{}{
				"type":        pType,
				"description": f.Description,
			}
		}
		schemas[s.Name] = map[string]interface{}{
			"type":        "object",
			"description": s.Description,
			"properties":  properties,
		}
	}

	// Map Routes to Paths
	for _, r := range doc.Routes {
		path := r.Path
		// Normalize parameters: convert /users/:id to /users/{id}
		parts := strings.Split(path, "/")
		for i, part := range parts {
			if strings.HasPrefix(part, ":") {
				parts[i] = "{" + strings.TrimPrefix(part, ":") + "}"
			}
		}
		normalizedPath := strings.Join(parts, "/")

		if _, exists := spec.Paths[normalizedPath]; !exists {
			spec.Paths[normalizedPath] = make(map[string]interface{})
		}

		pathItem := spec.Paths[normalizedPath].(map[string]interface{})
		method := strings.ToLower(r.Method)

		responses := map[string]interface{}{
			"200": map[string]interface{}{
				"description": "Success response",
			},
		}

		if r.OutputType != "" {
			responses["200"] = map[string]interface{}{
				"description": "Success response",
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"$ref": "#/components/schemas/" + r.OutputType,
						},
					},
				},
			}
		}

		operation := map[string]interface{}{
			"summary":     r.Description,
			"description": r.Description,
			"responses":   responses,
		}

		// If input type exists
		if r.InputType != "" && (method == "post" || method == "put" || method == "patch") {
			operation["requestBody"] = map[string]interface{}{
				"required": true,
				"content": map[string]interface{}{
					"application/json": map[string]interface{}{
						"schema": map[string]interface{}{
							"$ref": "#/components/schemas/" + r.InputType,
						},
					},
				},
			}
		}

		// Extract URL parameters
		parameters := make([]map[string]interface{}, 0)
		for _, part := range strings.Split(r.Path, "/") {
			if strings.HasPrefix(part, ":") {
				paramName := strings.TrimPrefix(part, ":")
				parameters = append(parameters, map[string]interface{}{
					"name":     paramName,
					"in":       "path",
					"required": true,
					"schema": map[string]interface{}{
						"type": "string",
					},
				})
			}
		}
		if len(parameters) > 0 {
			operation["parameters"] = parameters
		}

		pathItem[method] = operation
	}

	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(outputPath, data, 0644)
}
