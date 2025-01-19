## Status: alpha [![Go Reference](https://pkg.go.dev/badge/github.com/tassa-yoniso-manasi-karoto/go-aksharamukha.svg)](https://pkg.go.dev/github.com/tassa-yoniso-manasi-karoto/go-aksharamukha) 

Go bindings for [Aksharamukha](https://github.com/virtualvinodh/aksharamukha), a script converter and transliteration tool for various Indic and other scripts.

## Features

- Transliteration between 100+ scripts
- Romanization of text using ISO/academic standards
- Support for various Indic, Southeast Asian, and Middle Eastern scripts
- Docker-based deployment
- Script autoselection for language romanization
- Customizable transliteration options

**<p align="center"> ⚠️ While many scripts are supported I would not recommend using Aksharamukha for anything but romanization of indic languages or transliteration between indic languages. ⚠️ </p>**

## Quick Start

```go
import (
	"fmt"
	"log"
	ak "github.com/tassa-yoniso-manasi-karoto/go-aksharamukha"
)

func main() {
	// Initialize the service
	a, err := ak.NewAksharamukha()
	if err != nil {
		log.Fatal(err)
	}
	if err := a.Init(); err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	// Simple transliteration
	text := "नमस्ते"
	result, err := a.TransliterateSimple(text, ak.Devanagari, ak.Tamil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)

	// Romanization from a ISO-639 language code
	result, err = a.Romanize("नमस्ते", "hin")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
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
