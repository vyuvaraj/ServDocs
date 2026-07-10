// Package parser provides types and logic for parsing .srv source files
// into a structured SrvDoc representation for documentation generation.
package parser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// StructField represents a single field inside a .srv struct definition.
type StructField struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

// StructDef represents a struct definition parsed from a .srv file.
type StructDef struct {
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Fields      []StructField `json:"fields"`
}

// RouteDef represents an HTTP route definition parsed from a .srv file.
type RouteDef struct {
	Method      string   `json:"method"`
	Path        string   `json:"path"`
	InputType   string   `json:"input_type,omitempty"`
	OutputType  string   `json:"output_type,omitempty"`
	Description string   `json:"description,omitempty"`
	Middlewares []string `json:"middlewares,omitempty"`
}

// FnDef represents a function definition parsed from a .srv file.
type FnDef struct {
	Name        string `json:"name"`
	InputParams string `json:"input_params"`
	ReturnType  string `json:"return_type,omitempty"`
	Description string `json:"description,omitempty"`
}

// SrvDoc is the top-level document produced by parsing a .srv file.
// It aggregates all discovered structs, routes and functions.
type SrvDoc struct {
	Structs   []StructDef `json:"structs"`
	Routes    []RouteDef  `json:"routes"`
	Functions []FnDef     `json:"functions"`
}

var (
	structRegex = regexp.MustCompile(`^struct\s+(\w+)\s*\{`)
	fieldRegex  = regexp.MustCompile(`^\s*(\w+)\s*:\s*([^,;]+)`)
	routeRegex  = regexp.MustCompile(`^route\s+"?(GET|POST|PUT|DELETE)"?\s+"?([^\s"(]+)"?\s*\(([^)]*)\)(?:\s+use\s+\[([^\]]*)\])?(?:\s*->\s*(\S+))?`)
	fnRegex     = regexp.MustCompile(`^fn\s+(\w+)\s*\(([^)]*)\)\s*(?:->\s*(\S+))?`)
)

// ParseSrvFile opens and parses a .srv source file at filePath, returning
// a populated SrvDoc or an error if the file cannot be read.
func ParseSrvFile(filePath string) (*SrvDoc, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	doc := &SrvDoc{
		Structs:   make([]StructDef, 0),
		Routes:    make([]RouteDef, 0),
		Functions: make([]FnDef, 0),
	}

	scanner := bufio.NewScanner(file)
	var currentDocComments []string
	var activeStruct *StructDef

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "///") {
			comment := strings.TrimSpace(strings.TrimPrefix(line, "///"))
			currentDocComments = append(currentDocComments, comment)
			continue
		}

		if line == "" {
			continue
		}

		// Parse Struct
		if matches := structRegex.FindStringSubmatch(line); len(matches) > 1 {
			activeStruct = &StructDef{
				Name:        matches[1],
				Description: strings.Join(currentDocComments, " "),
				Fields:      make([]StructField, 0),
			}
			currentDocComments = nil
			continue
		}

		if activeStruct != nil {
			if strings.HasPrefix(line, "}") {
				doc.Structs = append(doc.Structs, *activeStruct)
				activeStruct = nil
			} else if matches := fieldRegex.FindStringSubmatch(line); len(matches) > 2 {
				desc := ""
				if len(currentDocComments) > 0 {
					desc = strings.Join(currentDocComments, " ")
					currentDocComments = nil
				}
				activeStruct.Fields = append(activeStruct.Fields, StructField{
					Name:        matches[1],
					Type:        strings.TrimSpace(matches[2]),
					Description: desc,
				})
			}
			continue
		}

		// Parse Route
		if matches := routeRegex.FindStringSubmatch(line); len(matches) > 2 {
			inputType := ""
			if matches[3] != "" {
				inputType = strings.TrimSpace(matches[3])
			}
			var middlewares []string
			if len(matches) > 4 && matches[4] != "" {
				parts := strings.Split(matches[4], ",")
				for _, p := range parts {
					trimmed := strings.TrimSpace(p)
					if trimmed != "" {
						middlewares = append(middlewares, trimmed)
					}
				}
			}
			outputType := ""
			if len(matches) > 5 {
				outputType = strings.TrimSpace(matches[5])
			}
			doc.Routes = append(doc.Routes, RouteDef{
				Method:      matches[1],
				Path:        matches[2],
				InputType:   inputType,
				OutputType:  outputType,
				Description: strings.Join(currentDocComments, " "),
				Middlewares: middlewares,
			})
			currentDocComments = nil
			continue
		}

		// Parse Function
		if matches := fnRegex.FindStringSubmatch(line); len(matches) > 2 {
			retType := ""
			if len(matches) > 3 {
				retType = strings.TrimSpace(matches[3])
			}
			doc.Functions = append(doc.Functions, FnDef{
				Name:        matches[1],
				InputParams: strings.TrimSpace(matches[2]),
				ReturnType:  retType,
				Description: strings.Join(currentDocComments, " "),
			})
			currentDocComments = nil
			continue
		}

		// Reset comments if line doesn't match anything
		currentDocComments = nil
	}

	return doc, scanner.Err()
}
