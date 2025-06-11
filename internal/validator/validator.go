package validator

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// SearchRequest represents a search request with validation
type SearchRequest struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

// Validator provides input validation functionality
type Validator struct {
	maxQueryLength int
	minQueryLength int
	maxLimit       int
	maxOffset      int
}

// NewValidator creates a new validator with default settings
func NewValidator() *Validator {
	return &Validator{
		maxQueryLength: 200,
		minQueryLength: 1,
		maxLimit:       100,
		maxOffset:      10000,
	}
}

// ValidateSearchRequest validates a search request
func (v *Validator) ValidateSearchRequest(req *SearchRequest) error {
	if err := v.validateQuery(req.Query); err != nil {
		return fmt.Errorf("invalid query: %w", err)
	}
	
	if err := v.validateLimit(req.Limit); err != nil {
		return fmt.Errorf("invalid limit: %w", err)
	}
	
	if err := v.validateOffset(req.Offset); err != nil {
		return fmt.Errorf("invalid offset: %w", err)
	}
	
	return nil
}

// validateQuery validates the search query
func (v *Validator) validateQuery(query string) error {
	// Trim whitespace
	query = strings.TrimSpace(query)
	
	// Check if empty
	if query == "" {
		return fmt.Errorf("query cannot be empty")
	}
	
	// Check length
	if utf8.RuneCountInString(query) < v.minQueryLength {
		return fmt.Errorf("query too short, minimum length is %d", v.minQueryLength)
	}
	
	if utf8.RuneCountInString(query) > v.maxQueryLength {
		return fmt.Errorf("query too long, maximum length is %d", v.maxQueryLength)
	}
	
	// Check for potentially dangerous patterns
	if err := v.checkDangerousPatterns(query); err != nil {
		return err
	}
	
	return nil
}

// validateLimit validates the limit parameter
func (v *Validator) validateLimit(limit int) error {
	if limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}
	
	if limit == 0 {
		// Set default limit
		limit = 10
	}
	
	if limit > v.maxLimit {
		return fmt.Errorf("limit too large, maximum is %d", v.maxLimit)
	}
	
	return nil
}

// validateOffset validates the offset parameter
func (v *Validator) validateOffset(offset int) error {
	if offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	
	if offset > v.maxOffset {
		return fmt.Errorf("offset too large, maximum is %d", v.maxOffset)
	}
	
	return nil
}

// checkDangerousPatterns checks for potentially dangerous input patterns
func (v *Validator) checkDangerousPatterns(query string) error {
	// List of dangerous patterns to check for
	dangerousPatterns := []string{
		`<script`,
		`javascript:`,
		`onload=`,
		`onerror=`,
		`eval\(`,
		`document\.`,
		`window\.`,
		`\.innerHTML`,
		`\.outerHTML`,
	}
	
	lowerQuery := strings.ToLower(query)
	
	for _, pattern := range dangerousPatterns {
		matched, err := regexp.MatchString(pattern, lowerQuery)
		if err != nil {
			continue // Skip invalid regex
		}
		if matched {
			return fmt.Errorf("query contains potentially dangerous content")
		}
	}
	
	return nil
}

// SanitizeQuery sanitizes the search query
func (v *Validator) SanitizeQuery(query string) string {
	// Trim whitespace
	query = strings.TrimSpace(query)
	
	// Remove control characters
	query = regexp.MustCompile(`[\x00-\x1f\x7f]`).ReplaceAllString(query, "")
	
	// Normalize whitespace
	query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")
	
	return query
}

// ValidateAndSanitizeQuery validates and sanitizes a search query
func (v *Validator) ValidateAndSanitizeQuery(query string) (string, error) {
	// First sanitize
	sanitized := v.SanitizeQuery(query)
	
	// Then validate
	if err := v.validateQuery(sanitized); err != nil {
		return "", err
	}
	
	return sanitized, nil
}

// SetLimits allows customizing validation limits
func (v *Validator) SetLimits(maxQueryLength, minQueryLength, maxLimit, maxOffset int) {
	if maxQueryLength > 0 {
		v.maxQueryLength = maxQueryLength
	}
	if minQueryLength > 0 {
		v.minQueryLength = minQueryLength
	}
	if maxLimit > 0 {
		v.maxLimit = maxLimit
	}
	if maxOffset >= 0 {
		v.maxOffset = maxOffset
	}
}

// Default validator instance
var defaultValidator = NewValidator()

// ValidateSearchRequest validates a search request using the default validator
func ValidateSearchRequest(req *SearchRequest) error {
	return defaultValidator.ValidateSearchRequest(req)
}

// SanitizeQuery sanitizes a query using the default validator
func SanitizeQuery(query string) string {
	return defaultValidator.SanitizeQuery(query)
}

// ValidateAndSanitizeQuery validates and sanitizes a query using the default validator
func ValidateAndSanitizeQuery(query string) (string, error) {
	return defaultValidator.ValidateAndSanitizeQuery(query)
}
