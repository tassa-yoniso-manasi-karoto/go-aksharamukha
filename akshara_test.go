package aksharamukha

import (
	"testing"
	"time"
	"strings"
)

func TestRomanization(t *testing.T) {
	// Initialize Aksharamukha
	a, err := NewAksharamukha()
	if err != nil {
		t.Fatalf("Failed to create Aksharamukha instance: %v", err)
	}
	if err := a.Init(); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	//defer a.Close()

	// Wait for service to be ready
	time.Sleep(5 * time.Second)

	tests := []struct {
		name     string
		text     string
		lang     string
		expected string
	}{
		{
			name:     "Hindi",
			text:     "नमस्ते",
			lang:     "hin",
			expected: "namastē",
		},
		{
			name:     "Bengali",
			text:     "নমস্কার",
			lang:     "ben",
			expected: "namaskāra",
		},
		{
			name:     "Arabic",
			text:     "السَّلامُ عَلَيْكُمْ",
			lang:     "ara",
			expected: "ʾls꞉alʾmu ʿalaŷkum",
		},
		{
			name:     "Sanskrit (Devanagari)",
			text:     "संस्कृतम्",
			lang:     "san",
			expected: "saṁskr̥tam",
		},
		{
			name:     "Tamil",
			text:     "வணக்கம்",
			lang:     "tam",
			expected: "vaṇakkam",
		},
		{
			name:     "Telugu",
			text:     "నమస్కారం",
			lang:     "tel",
			expected: "namaskāraṁ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := a.Romanize(tt.text, tt.lang)
			if err != nil {
				t.Errorf("Romanize() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Romanize() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransliteration(t *testing.T) {
	// Initialize Aksharamukha
	a, err := NewAksharamukha()
	if err != nil {
		t.Fatalf("Failed to create Aksharamukha instance: %v", err)
	}
	if err := a.Init(); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	//defer a.Close()

	// Wait for service to be ready
	time.Sleep(5 * time.Second)

	tests := []struct {
		name     string
		text     string
		from     Script
		to       Script
		expected string
	}{
		{
			name:     "Devanagari to Tamil",
			text:     "नमस्ते",
			from:     Devanagari,
			to:       Tamil,
			expected: "நமஸ்தே",
		},
		{
			name:     "Bengali to Devanagari",
			text:     "নমস্কার",
			from:     Bengali,
			to:       Devanagari,
			expected: "नमस्कार",
		},
		{
			name:     "IAST to Devanagari",
			text:     "namaste",
			from:     IAST,
			to:       Devanagari,
			expected: "नमस्ते",
		},
		{
			name:     "Tamil to Telugu",
			text:     "வணக்கம்",
			from:     Tamil,
			to:       Telugu,
			expected: "వణక్కమ్",
		},
		{
			name:     "Devanagari to Malayalam",
			text:     "नमस्कार",
			from:     Devanagari,
			to:       Malayalam,
			expected: "നമസ്കാര",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := a.TransliterateSimple(tt.text, tt.from, tt.to)
			if err != nil {
				t.Errorf("TransliterateSimple() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("TransliterateSimple() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransliterationWithOptions(t *testing.T) {
	// Initialize Aksharamukha
	a, err := NewAksharamukha()
	if err != nil {
		t.Fatalf("Failed to create Aksharamukha instance: %v", err)
	}
	if err := a.Init(); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	//defer a.Close()

	// Wait for service to be ready
	time.Sleep(5 * time.Second)

	tests := []struct {
		name     string
		text     string
		from     Script
		to       Script
		opts     TransliterationOptions
		expected string
	}{
		{
			name: "Sanskrit to Telugu without nativization",
			text: "भगवद्गीता",
			from: Devanagari,
			to:   Telugu,
			opts: TransliterationOptions{
				Nativize: false,
			},
			expected: "భగవద్గీతా",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := a.Transliterate(tt.text, tt.from, tt.to, tt.opts)
			if err != nil {
				t.Errorf("Transliterate() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Transliterate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInvalidInputs(t *testing.T) {
	a, err := NewAksharamukha()
	if err != nil {
		t.Fatalf("Failed to create Aksharamukha instance: %v", err)
	}
	if err := a.Init(); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	//defer a.Close()

	// Invalid language code
	_, err = a.Romanize("test", "invalid")
	if err == nil {
		t.Error("Expected error for invalid language code, got nil")
	}

	// Empty text
	_, err = a.Romanize("", "hin")
	if err == nil {
		t.Error("Expected error for empty text, got nil")
	}

	// Invalid script combination
	_, err = a.TransliterateSimple("test", Script("InvalidScript"), Devanagari)
	if err == nil {
		t.Error("Expected error for invalid script, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid source script") {
		t.Errorf("Expected 'invalid source script' error, got: %v", err)
	}

	// Invalid target script
	_, err = a.TransliterateSimple("test", Devanagari, Script("InvalidScript"))
	if err == nil {
		t.Error("Expected error for invalid target script, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid target script") {
		t.Errorf("Expected 'invalid target script' error, got: %v", err)
	}
}
