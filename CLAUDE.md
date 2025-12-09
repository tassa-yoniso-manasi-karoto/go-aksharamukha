# Go-Aksharamukha Development Guide

## Build & Test Commands
- Run tests: `go test ./...`
- Run specific test: `go test -run TestName`
- Start containers only: `docker compose up`
- Force container rebuild: Use `InitRecreate(ctx, true)` in code
- Format code: `gofmt -w *.go`
- Lint code: `go vet ./...`

## Code Style Guidelines
- Format with gofmt
- Errors: Use `fmt.Errorf("message: %w", err)` for error wrapping
- Logging: Use zerolog package (`github.com/rs/zerolog`)
- Variable naming: camelCase for private, PascalCase for exported
- Struct fields alignment: Align adjacent field names and tags
- Error handling: Check all errors, don't use panic
- Documentation: Add doc comments for all exported functions/types
- Imports: Group standard library, third-party, and local imports
- Testing: Write tests for public APIs
- Constants: Use package-level const/var blocks for related values

## Modern API Usage (Multiple Instances Support)
```go
// Create a new Aksharamukha manager with custom settings
ctx := context.Background()
manager, err := aksharamukha.NewManager(ctx, 
    aksharamukha.WithProjectName("aksharamukha-custom"),
    aksharamukha.WithQueryTimeout(10 * time.Minute))
if err != nil {
    log.Fatal(err)
}

// Initialize the Docker container
if err := manager.Init(ctx); err != nil {
    log.Fatal(err)
}

// Transliterate text using the manager
result, err := manager.Translit(ctx, "नमस्ते", aksharamukha.Devanagari, aksharamukha.Tamil)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)

// Clean up when done
defer manager.Close()
```

## Context Guidelines
- Always pass context as first parameter to functions
- Use timeout-wrapped contexts for HTTP API calls
- Never store context in struct fields
- For long-running operations, derive child contexts with appropriate timeouts

## Script Support
- Use the predefined Script constants (e.g., `aksharamukha.Devanagari`)
- For language-based transliteration, use the ISO 639 codes with `Roman()` function
- Check script validity with `IsValidScript()` before making API calls
- See `scripts.go` for all supported scripts and their romanization mappings