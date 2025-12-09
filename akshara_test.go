package aksharamukha

import (
	"context"
	"testing"
	"time"
	"strings"
)

func TestRomanizationBackwardCompatible(t *testing.T) {
	t.Skip("Skipping test that requires Docker container - run manually")
	
	// Initialize Aksharamukha using backward compatible API
	if err := Init(); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	defer Close()

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
			name:     "Sanskrit (Devanagari)",
			text:     "संस्कृतम्",
			lang:     "san",
			expected: "saṁskr̥tam",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Roman(tt.text, tt.lang)
			if err != nil {
				t.Errorf("Roman() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Roman() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRomanizationWithContext(t *testing.T) {
	t.Skip("Skipping test that requires Docker container - run manually")
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	
	// Initialize Aksharamukha using context-aware API
	if err := InitWithContext(ctx); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	defer Close()

	tests := []struct {
		name     string
		text     string
		lang     string
		expected string
	}{
		{
			name:     "Bengali",
			text:     "নমস্কার",
			lang:     "ben",
			expected: "namaskāra",
		},
		{
			name:     "Tamil",
			text:     "வணக்கம்",
			lang:     "tam",
			expected: "vaṇakkam",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := RomanWithContext(ctx, tt.text, tt.lang, DefaultOptions())
			if err != nil {
				t.Errorf("RomanWithContext() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("RomanWithContext() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransliterationBackwardCompatible(t *testing.T) {
	t.Skip("Skipping test that requires Docker container - run manually")
	
	// Initialize Aksharamukha using backward compatible API
	if err := Init(); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	defer Close()

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
			name:     "IAST to Devanagari",
			text:     "namaste",
			from:     IAST,
			to:       Devanagari,
			expected: "नमस्ते",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Translit(tt.text, tt.from, tt.to)
			if err != nil {
				t.Errorf("Translit() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Translit() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTransliterationWithContext(t *testing.T) {
	t.Skip("Skipping test that requires Docker container - run manually")
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	
	// Initialize Aksharamukha using context-aware API
	if err := InitWithContext(ctx); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	defer Close()

	tests := []struct {
		name     string
		text     string
		from     Script
		to       Script
		expected string
	}{
		{
			name:     "Bengali to Devanagari",
			text:     "নমস্কার",
			from:     Bengali,
			to:       Devanagari,
			expected: "नमस्कार",
		},
		{
			name:     "Tamil to Telugu",
			text:     "வணக்கம்",
			from:     Tamil,
			to:       Telugu,
			expected: "వణక్కమ్",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TranslitWithContext(ctx, tt.text, tt.from, tt.to, DefaultOptions())
			if err != nil {
				t.Errorf("TranslitWithContext() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("TranslitWithContext() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestManagerAPI(t *testing.T) {
	t.Skip("Skipping test that requires Docker container - run manually")
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	
	// Create a custom manager
	manager, err := NewManager(ctx, 
		WithProjectName("aksharamukha-test"),
		WithQueryTimeout(30*time.Second))
	if err != nil {
		t.Fatalf("Failed to create Aksharamukha manager: %v", err)
	}
	
	// Initialize manager
	if err := manager.Init(ctx); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha manager: %v", err)
	}
	defer manager.Close()

	// Test transliteration with manager
	text := "नमस्ते"
	result, err := manager.Translit(ctx, text, Devanagari, Tamil, DefaultOptions())
	if err != nil {
		t.Errorf("Manager.Translit() error = %v", err)
		return
	}
	
	expected := "நமஸ்தே"
	if result != expected {
		t.Errorf("Manager.Translit() = %v, want %v", result, expected)
	}
}

func TestMultipleManagers(t *testing.T) {
	t.Skip("Skipping test that requires Docker container - run manually")
	
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	
	// Create first manager
	manager1, err := NewManager(ctx, 
		WithProjectName("aksharamukha-test1"))
	if err != nil {
		t.Fatalf("Failed to create first Aksharamukha manager: %v", err)
	}
	
	// Initialize first manager
	if err := manager1.Init(ctx); err != nil {
		t.Fatalf("Failed to initialize first Aksharamukha manager: %v", err)
	}
	defer manager1.Close()
	
	// Create second manager
	manager2, err := NewManager(ctx, 
		WithProjectName("aksharamukha-test2"))
	if err != nil {
		t.Fatalf("Failed to create second Aksharamukha manager: %v", err)
	}
	
	// Initialize second manager
	if err := manager2.Init(ctx); err != nil {
		t.Fatalf("Failed to initialize second Aksharamukha manager: %v", err)
	}
	defer manager2.Close()

	// Test both managers concurrently
	text1 := "नमस्ते"
	expected1 := "நமஸ்தே"
	
	text2 := "संस्कृतम्"
	expected2 := "සංස්කෘතම්"
	
	// Use first manager
	result1, err := manager1.Translit(ctx, text1, Devanagari, Tamil, DefaultOptions())
	if err != nil {
		t.Errorf("Manager1.Translit() error = %v", err)
	} else if result1 != expected1 {
		t.Errorf("Manager1.Translit() = %v, want %v", result1, expected1)
	}
	
	// Use second manager
	result2, err := manager2.Translit(ctx, text2, Devanagari, Sinhala, DefaultOptions())
	if err != nil {
		t.Errorf("Manager2.Translit() error = %v", err)
	} else if result2 != expected2 {
		t.Errorf("Manager2.Translit() = %v, want %v", result2, expected2)
	}
}

func TestInvalidInputs(t *testing.T) {
	t.Skip("Skipping test that requires Docker container - run manually")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := InitWithContext(ctx); err != nil {
		t.Fatalf("Failed to initialize Aksharamukha: %v", err)
	}
	defer Close()

	// Invalid language code
	_, err := RomanWithContext(ctx, "test", "invalid", DefaultOptions())
	if err == nil {
		t.Error("Expected error for invalid language code, got nil")
	}

	// Empty text
	_, err = RomanWithContext(ctx, "", "hin", DefaultOptions())
	if err == nil {
		t.Error("Expected error for empty text, got nil")
	}

	// Invalid script combination
	_, err = TranslitWithContext(ctx, "test", Script("InvalidScript"), Devanagari, DefaultOptions())
	if err == nil {
		t.Error("Expected error for invalid script, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid source script") {
		t.Errorf("Expected 'invalid source script' error, got: %v", err)
	}

	// Invalid target script
	_, err = TranslitWithContext(ctx, "test", Devanagari, Script("InvalidScript"), DefaultOptions())
	if err == nil {
		t.Error("Expected error for invalid target script, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "invalid target script") {
		t.Errorf("Expected 'invalid target script' error, got: %v", err)
	}
	
	// Context cancellation
	cancelCtx, cancelFunc := context.WithCancel(context.Background())
	cancelFunc() // Cancel immediately
	_, err = RomanWithContext(cancelCtx, "test", "hin", DefaultOptions())
	if err == nil || !strings.Contains(err.Error(), "context canceled") {
		t.Errorf("Expected context cancellation error, got: %v", err)
	}
}