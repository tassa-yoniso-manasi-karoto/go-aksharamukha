## Status: alpha [![Go Reference](https://pkg.go.dev/badge/github.com/tassa-yoniso-manasi-karoto/go-aksharamukha.svg)](https://pkg.go.dev/github.com/tassa-yoniso-manasi-karoto/go-aksharamukha) 

Go bindings for [Aksharamukha](https://github.com/virtualvinodh/aksharamukha), a script converter and transliteration tool for various Indic and other scripts.

## Features

- Transliteration between 100+ scripts
- Romanization of text using ISO/academic standards
- Support for various Indic, Southeast Asian, and Middle Eastern scripts
- Docker-based deployment
- Script autoselection for language romanization
- Customizable transliteration options
- Context support for proper cancellation and timeouts
- Multiple instance support for concurrent processing

**<p align="center"> ⚠️ While many scripts are supported I would not recommend using Aksharamukha for anything but romanization of indic languages or transliteration between indic languages. ⚠️ </p>**

## Quick Start

### Simple Usage (Backward Compatible)

```go
import (
	"fmt"
	"log"
	ak "github.com/tassa-yoniso-manasi-karoto/go-aksharamukha"
)

func main() {
	// Initialize the environment (downloads, builds and starts containers)
	if err := ak.Init(); err != nil {
		log.Fatal(err)
	}
	defer ak.Close()

	// Simple transliteration
	text := "नमस्ते"
	result, err := ak.Translit(text, ak.Devanagari, ak.Tamil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	// Romanization from a ISO-639 language code
	// Uses the default romanization scheme per language
	result, err = ak.Roman("नमस्ते", "hin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
```

### Context-Aware API

```go
import (
	"context"
	"fmt"
	"log"
	"time"
	
	ak "github.com/tassa-yoniso-manasi-karoto/go-aksharamukha"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	// Initialize with context
	if err := ak.InitWithContext(ctx); err != nil {
		log.Fatal(err)
	}
	defer ak.Close()

	// Transliteration with context
	text := "नमस्ते"
	result, err := ak.TranslitWithContext(ctx, text, ak.Devanagari, ak.Tamil, ak.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	
	// Romanization with context
	result, err = ak.RomanWithContext(ctx, "नमस्ते", "hin", ak.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}
```

### Manager-Based API (Multiple Instances)

```go
import (
	"context"
	"fmt"
	"log"
	"time"
	
	ak "github.com/tassa-yoniso-manasi-karoto/go-aksharamukha"
)

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	
	// Create a custom manager with options
	manager, err := ak.NewManager(ctx, 
		ak.WithProjectName("aksharamukha-custom"),
		ak.WithQueryTimeout(30*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	
	// Initialize the environment
	if err := manager.Init(ctx); err != nil {
		log.Fatal(err)
	}
	defer manager.Close()

	// Transliterate using the manager
	text := "नमस्ते"
	result, err := manager.Translit(ctx, text, ak.Devanagari, ak.Tamil, ak.DefaultOptions())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	
	// Create a second manager instance if needed
	manager2, err := ak.NewManager(ctx, 
		ak.WithProjectName("aksharamukha-second"))
	if err != nil {
		log.Fatal(err)
	}
	
	if err := manager2.Init(ctx); err != nil {
		log.Fatal(err)
	}
	defer manager2.Close()
	
	// Now you can use both managers concurrently for separate tasks
}
```

### Output

```
நமஸ்தே
namastē
```

## Requirements

- Go 1.16 or later
- **Installed Docker Engine (linux) or Docker Desktop (windows/mac)**
- Internet connection (for initial container pull)
