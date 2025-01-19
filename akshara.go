
package aksharamukha

import (
	"io"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	
	"github.com/gookit/color"
	"github.com/k0kubun/pp"
	iso "github.com/barbashov/iso639-3"
)

// TransliterationOptions holds configuration for the transliteration process
type TransliterationOptions struct {
	// If false, prevents automatic nativization according to output script conventions
	Nativize bool
	// Options applied before transliteration
	PreOptions []string
	// Options applied after transliteration
	PostOptions []string
}

// DefaultOptions returns the default transliteration options
func DefaultOptions() TransliterationOptions {
	return TransliterationOptions{}
		//Nativize: true,
	//}
}

// Transliterate converts text from one script to another
// Transliterate converts text from one script to another
func (a *Aksharamukha) Transliterate(text string, from, to Script, opts TransliterationOptions) (string, error) {
	if text == "" {
		return "", fmt.Errorf("empty text provided")
	}

	// Validate scripts if provided
	if from != "" && !IsValidScript(from) {
		return "", fmt.Errorf("invalid source script: %s", from)
	}
	if !IsValidScript(to) {
		return "", fmt.Errorf("invalid target script: %s", to)
	}


	// Build the query URL
	baseURL := "http://localhost:8085/api/public"
	params := url.Values{}
	
	// Required parameters
	params.Set("text", text)
	params.Set("target", string(to))
	
	// Optional source script (if not provided, system will auto-detect)
	if from != "" {
		params.Set("source", string(from))
	}
	
	// Optional nativization parameter (only set if false, as true is default)
	if !opts.Nativize {
		params.Set("nativize", "false")
	}
	
	// Optional pre-options
	if len(opts.PreOptions) > 0 {
		params.Set("preoptions", strings.Join(opts.PreOptions, ","))
	}
	
	// Optional post-options
	if len(opts.PostOptions) > 0 {
		params.Set("postoptions", strings.Join(opts.PostOptions, ","))
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	resp, err := client.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Return the response as is, since it's plain text
	result := string(body)
	result = strings.TrimSpace(result) // Remove any leading/trailing whitespace

	if result == "" {
		return "", fmt.Errorf("empty response received")
	}

	return result, nil
}

// TransliterateSimple is a convenience function for simple transliteration without options
func (a *Aksharamukha) TransliterateSimple(text string, from, to Script) (string, error) {
	return a.Transliterate(text, from, to, DefaultOptions())
}


// Romanize converts text from a given language to its romanized form
func (a *Aksharamukha) Romanize(text, languageCode string) (string, error) {
	// Validate and standardize the language code
	stdLang, ok := IsValidISO639(languageCode)
	if !ok {
		return "", fmt.Errorf("\"%s\" isn't a ISO-639 language code", languageCode)
	}

	// Get the script for the language
	scripts, exists := Lang2Scripts[stdLang]
	if !exists {
		return "", fmt.Errorf("no script mapping found for language code %s", stdLang)
	}
	if len(scripts) == 0 {
		return "", fmt.Errorf("empty script list for language code %s", stdLang)
	}

	// Get the primary script (first in the list)
	sourceScript := Script(scripts[0])

	// Get the romanization scheme for the script
	romanScheme, exists := Script2RomanScheme[string(sourceScript)]
	if !exists {
		return "", fmt.Errorf("no romanization scheme found for script %s", sourceScript)
	}

	// Use the existing Transliterate function with default options
	result, err := a.Transliterate(text, sourceScript, Script(romanScheme), DefaultOptions())
	if err != nil {
		return "", fmt.Errorf("romanization failed: %w", err)
	}

	return result, nil
}

// RomanizeWithOptions is like Romanize but allows customization of the transliteration options
func (a *Aksharamukha) RomanizeWithOptions(text, languageCode string, opts TransliterationOptions) (string, error) {
	// Validate and standardize the language code
	stdLang, ok := IsValidISO639(languageCode)
	if !ok {
		return "", fmt.Errorf("\"%s\" isn't a ISO-639 language code", languageCode)
	}

	// Get the script for the language
	scripts, exists := Lang2Scripts[stdLang]
	if !exists {
		return "", fmt.Errorf("no script mapping found for language code %s", stdLang)
	}
	if len(scripts) == 0 {
		return "", fmt.Errorf("empty script list for language code %s", stdLang)
	}

	// Get the primary script (first in the list)
	sourceScript := Script(scripts[0])

	// Get the romanization scheme for the script
	romanScheme, exists := Script2RomanScheme[string(sourceScript)]
	if !exists {
		return "", fmt.Errorf("no romanization scheme found for script %s", sourceScript)
	}

	// Use the existing Transliterate function with provided options
	result, err := a.Transliterate(text, sourceScript, Script(romanScheme), opts)
	if err != nil {
		return "", fmt.Errorf("romanization failed: %w", err)
	}

	return result, nil
}

func IsValidISO639(lang string) (stdLang string, ok bool) {
	code := iso.FromAnyCode(lang)
	if code == nil {
		return
	}
	return code.Part3, true
}


func placeholder() {
	color.Redln(" 𝒻*** 𝓎ℴ𝓊 𝒸ℴ𝓂𝓅𝒾𝓁ℯ𝓇")
	pp.Println("𝓯*** 𝔂𝓸𝓾 𝓬𝓸𝓶𝓹𝓲𝓵𝓮𝓻")
}
