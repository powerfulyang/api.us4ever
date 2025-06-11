package search

import (
	"context"
	"time"
)

// SearchService defines the interface for search operations
type SearchService interface {
	// SearchKeeps searches for keeps with the given query and options
	SearchKeeps(ctx context.Context, query string, opts ...SearchOption) (*SearchResult, error)
	
	// SearchMoments searches for moments with the given query and options
	SearchMoments(ctx context.Context, query string, opts ...SearchOption) (*SearchResult, error)
	
	// HealthCheck checks if the search service is healthy
	HealthCheck(ctx context.Context) error
	
	// Close closes the search service and cleans up resources
	Close() error
}

// SearchResult represents the result of a search operation
type SearchResult struct {
	// Query is the original search query
	Query string `json:"query"`
	
	// Total is the total number of results found
	Total int `json:"total"`
	
	// Hits contains the actual search results
	Hits []SearchHit `json:"hits"`
	
	// Duration is how long the search took
	Duration time.Duration `json:"duration"`
	
	// Aggregations contains any aggregation results
	Aggregations map[string]interface{} `json:"aggregations,omitempty"`
}

// SearchHit represents a single search result
type SearchHit struct {
	// ID is the unique identifier of the document
	ID string `json:"id"`
	
	// Score is the relevance score
	Score float64 `json:"score"`
	
	// Source contains the document data
	Source map[string]interface{} `json:"source"`
	
	// Highlight contains highlighted text snippets
	Highlight map[string][]string `json:"highlight,omitempty"`
	
	// Index is the name of the index this result came from
	Index string `json:"index"`
}

// SearchConfig holds the configuration for a search operation
type SearchConfig struct {
	// Limit is the maximum number of results to return
	Limit int
	
	// Offset is the number of results to skip
	Offset int
	
	// SortBy specifies the field to sort by
	SortBy string
	
	// SortOrder specifies the sort order (asc/desc)
	SortOrder string
	
	// Fields specifies which fields to include in the response
	Fields []string
	
	// Filters contains additional filters to apply
	Filters map[string]interface{}
	
	// IncludeHighlight specifies whether to include highlighting
	IncludeHighlight bool
	
	// IncludeAggregations specifies whether to include aggregations
	IncludeAggregations bool
	
	// Timeout specifies the maximum time to wait for results
	Timeout time.Duration
}

// SearchOption is a function that modifies SearchConfig
type SearchOption func(*SearchConfig)

// WithLimit sets the maximum number of results to return
func WithLimit(limit int) SearchOption {
	return func(cfg *SearchConfig) {
		cfg.Limit = limit
	}
}

// WithOffset sets the number of results to skip
func WithOffset(offset int) SearchOption {
	return func(cfg *SearchConfig) {
		cfg.Offset = offset
	}
}

// WithSort sets the sort field and order
func WithSort(field, order string) SearchOption {
	return func(cfg *SearchConfig) {
		cfg.SortBy = field
		cfg.SortOrder = order
	}
}

// WithFields sets the fields to include in the response
func WithFields(fields ...string) SearchOption {
	return func(cfg *SearchConfig) {
		cfg.Fields = fields
	}
}

// WithFilter adds a filter to the search
func WithFilter(key string, value interface{}) SearchOption {
	return func(cfg *SearchConfig) {
		if cfg.Filters == nil {
			cfg.Filters = make(map[string]interface{})
		}
		cfg.Filters[key] = value
	}
}

// WithHighlight enables highlighting in search results
func WithHighlight() SearchOption {
	return func(cfg *SearchConfig) {
		cfg.IncludeHighlight = true
	}
}

// WithAggregations enables aggregations in search results
func WithAggregations() SearchOption {
	return func(cfg *SearchConfig) {
		cfg.IncludeAggregations = true
	}
}

// WithTimeout sets the search timeout
func WithTimeout(timeout time.Duration) SearchOption {
	return func(cfg *SearchConfig) {
		cfg.Timeout = timeout
	}
}

// DefaultSearchConfig returns a SearchConfig with sensible defaults
func DefaultSearchConfig() *SearchConfig {
	return &SearchConfig{
		Limit:               10,
		Offset:              0,
		SortBy:              "_score",
		SortOrder:           "desc",
		Fields:              nil, // Include all fields
		Filters:             make(map[string]interface{}),
		IncludeHighlight:    true,
		IncludeAggregations: false,
		Timeout:             30 * time.Second,
	}
}

// ApplyOptions applies the given options to the config
func (cfg *SearchConfig) ApplyOptions(opts ...SearchOption) {
	for _, opt := range opts {
		opt(cfg)
	}
}

// Validate validates the search configuration
func (cfg *SearchConfig) Validate() error {
	if cfg.Limit < 0 {
		cfg.Limit = 10
	}
	if cfg.Limit > 1000 {
		cfg.Limit = 1000
	}
	
	if cfg.Offset < 0 {
		cfg.Offset = 0
	}
	
	if cfg.SortOrder != "asc" && cfg.SortOrder != "desc" {
		cfg.SortOrder = "desc"
	}
	
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	
	return nil
}

// SearchError represents an error that occurred during search
type SearchError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Query   string `json:"query,omitempty"`
}

// Error implements the error interface
func (e *SearchError) Error() string {
	return e.Message
}

// NewSearchError creates a new search error
func NewSearchError(errorType, message, query string, code int) *SearchError {
	return &SearchError{
		Type:    errorType,
		Message: message,
		Code:    code,
		Query:   query,
	}
}

// Common search errors
var (
	ErrInvalidQuery     = NewSearchError("InvalidQuery", "Invalid search query", "", 400)
	ErrIndexNotFound    = NewSearchError("IndexNotFound", "Search index not found", "", 404)
	ErrTimeout          = NewSearchError("Timeout", "Search request timeout", "", 408)
	ErrServiceUnavailable = NewSearchError("ServiceUnavailable", "Search service unavailable", "", 503)
)
