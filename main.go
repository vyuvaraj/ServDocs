package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"servdocs/pkg/generator"
	"servdocs/pkg/openapi"
	"servdocs/pkg/parser"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "generate":
		generateCmd := flag.NewFlagSet("generate", flag.ExitOnError)
		input := generateCmd.String("input", "example.srv", "Path to input .srv file")
		output := generateCmd.String("output", "docs.html", "Path to output HTML file")
		title := generateCmd.String("title", "Servverse", "Documentation Title")
		versionTag := generateCmd.String("version-tag", "", "Version tag (e.g. v1.0.0)")
		outDir := generateCmd.String("out-dir", "", "Output directory for versioned docs")
		_ = generateCmd.Parse(os.Args[2:])

		doc, err := parser.ParseSrvFile(*input)
		if err != nil {
			log.Fatalf("failed to parse file: %v", err)
		}

		if *outDir != "" && *versionTag != "" {
			err := os.MkdirAll(filepath.Join(*outDir, *versionTag), 0755)
			if err != nil {
				log.Fatalf("failed to create directory: %v", err)
			}
			*output = filepath.Join(*outDir, *versionTag, "index.html")
		}

		if err := generator.GenerateHtml(doc, *title, *output, *outDir, *versionTag); err != nil {
			log.Fatalf("failed to generate HTML: %v", err)
		}
		fmt.Printf("Documentation site generated successfully at %s\n", *output)

	case "openapi":
		openapiCmd := flag.NewFlagSet("openapi", flag.ExitOnError)
		input := openapiCmd.String("input", "example.srv", "Path to input .srv file")
		output := openapiCmd.String("output", "openapi.json", "Path to output OpenAPI JSON file")
		title := openapiCmd.String("title", "Servverse API", "API Title")
		_ = openapiCmd.Parse(os.Args[2:])

		doc, err := parser.ParseSrvFile(*input)
		if err != nil {
			log.Fatalf("failed to parse file: %v", err)
		}

		if err := openapi.Generate(doc, *title, *output); err != nil {
			log.Fatalf("failed to generate OpenAPI spec: %v", err)
		}
		fmt.Printf("OpenAPI specification generated successfully at %s\n", *output)

	case "serve":
		serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
		input := serveCmd.String("input", "example.srv", "Path to input .srv file")
		port := serveCmd.String("port", "3000", "Local server port")
		title := serveCmd.String("title", "Servverse", "Documentation Title")
		_ = serveCmd.Parse(os.Args[2:])

		doc, err := parser.ParseSrvFile(*input)
		if err != nil {
			log.Fatalf("failed to parse file: %v", err)
		}

		tempFile := "index_temp.html"
		if err := generator.GenerateHtml(doc, *title, tempFile, "", ""); err != nil {
			log.Fatalf("failed to generate HTML: %v", err)
		}
		defer os.Remove(tempFile)

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, tempFile)
		})

		fmt.Printf("Serving documentation at http://localhost:%s\n", *port)
		if err := http.ListenAndServe(":"+*port, nil); err != nil {
			log.Fatalf("failed to serve documentation: %v", err)
		}

	case "client":
		clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
		input := clientCmd.String("input", "example.srv", "Path to input .srv file")
		lang := clientCmd.String("lang", "typescript", "Target SDK language (typescript|dart|swift)")
		output := clientCmd.String("output", "sdk", "Output directory for the SDK")
		_ = clientCmd.Parse(os.Args[2:])

		doc, err := parser.ParseSrvFile(*input)
		if err != nil {
			log.Fatalf("failed to parse file: %v", err)
		}

		if err := generator.GenerateClientSDK(doc, *lang, *output); err != nil {
			log.Fatalf("failed to generate client SDK: %v", err)
		}
		fmt.Printf("Client SDK generated successfully at %s\n", *output)

	default:
		fmt.Printf("unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Usage: servdocs <command> [arguments]")
	fmt.Println("Commands:")
	fmt.Println("  generate   Generate static HTML documentation")
	fmt.Println("  openapi    Generate OpenAPI 3.0 specification JSON")
	fmt.Println("  serve      Serve the documentation site locally")
	fmt.Println("  client     Generate TypeScript, Dart, or Swift SDK clients")
}
