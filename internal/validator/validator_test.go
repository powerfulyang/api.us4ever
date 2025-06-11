package validator

import (
	"testing"
)

func TestValidateSearchRequest(t *testing.T) {
	validator := NewValidator()
	
	tests := []struct {
		name    string
		request SearchRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: SearchRequest{
				Query:  "test query",
				Limit:  10,
				Offset: 0,
			},
			wantErr: false,
		},
		{
			name: "empty query",
			request: SearchRequest{
				Query:  "",
				Limit:  10,
				Offset: 0,
			},
			wantErr: true,
		},
		{
			name: "query too long",
			request: SearchRequest{
				Query:  string(make([]byte, 300)), // Longer than max
				Limit:  10,
				Offset: 0,
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			request: SearchRequest{
				Query:  "test",
				Limit:  -1,
				Offset: 0,
			},
			wantErr: true,
		},
		{
			name: "limit too large",
			request: SearchRequest{
				Query:  "test",
				Limit:  1000,
				Offset: 0,
			},
			wantErr: true,
		},
		{
			name: "negative offset",
			request: SearchRequest{
				Query:  "test",
				Limit:  10,
				Offset: -1,
			},
			wantErr: true,
		},
		{
			name: "dangerous script tag",
			request: SearchRequest{
				Query:  "<script>alert('xss')</script>",
				Limit:  10,
				Offset: 0,
			},
			wantErr: true,
		},
		{
			name: "javascript injection",
			request: SearchRequest{
				Query:  "javascript:alert('xss')",
				Limit:  10,
				Offset: 0,
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateSearchRequest(&tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSearchRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeQuery(t *testing.T) {
	validator := NewValidator()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal query",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "query with extra spaces",
			input:    "  hello    world  ",
			expected: "hello world",
		},
		{
			name:     "query with control characters",
			input:    "hello\x00\x1fworld",
			expected: "helloworld",
		},
		{
			name:     "query with tabs and newlines",
			input:    "hello\t\nworld",
			expected: "hello world",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.SanitizeQuery(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeQuery() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestValidateAndSanitizeQuery(t *testing.T) {
	validator := NewValidator()
	
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "valid query with spaces",
			input:    "  hello world  ",
			expected: "hello world",
			wantErr:  false,
		},
		{
			name:     "empty query after sanitization",
			input:    "   ",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "dangerous content",
			input:    "<script>alert('xss')</script>",
			expected: "",
			wantErr:  true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validator.ValidateAndSanitizeQuery(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAndSanitizeQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("ValidateAndSanitizeQuery() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSetLimits(t *testing.T) {
	validator := NewValidator()
	
	// Test setting custom limits
	validator.SetLimits(500, 2, 50, 5000)
	
	// Test with query that would be valid with new limits
	req := &SearchRequest{
		Query:  "ab", // 2 characters (new minimum)
		Limit:  50,   // New maximum
		Offset: 5000, // New maximum
	}
	
	err := validator.ValidateSearchRequest(req)
	if err != nil {
		t.Errorf("ValidateSearchRequest() with custom limits failed: %v", err)
	}
	
	// Test with query that exceeds new limits
	req2 := &SearchRequest{
		Query:  "a", // Below new minimum
		Limit:  51,  // Above new maximum
		Offset: 5001, // Above new maximum
	}
	
	err2 := validator.ValidateSearchRequest(req2)
	if err2 == nil {
		t.Error("ValidateSearchRequest() should have failed with custom limits")
	}
}

func BenchmarkValidateSearchRequest(b *testing.B) {
	validator := NewValidator()
	req := &SearchRequest{
		Query:  "test query for benchmarking",
		Limit:  10,
		Offset: 0,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.ValidateSearchRequest(req)
	}
}

func BenchmarkSanitizeQuery(b *testing.B) {
	validator := NewValidator()
	query := "  hello    world with   extra   spaces  "
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.SanitizeQuery(query)
	}
}
