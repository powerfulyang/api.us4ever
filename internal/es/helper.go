package es

import (
	"bytes"
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type SearchParams struct {
	Keyword string
	Fields  []string
	Index   string
}

func BuildBody(p SearchParams) *bytes.Buffer {
	requireMatch := false
	req := search.Request{
		Query: &types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:  p.Keyword,
				Fields: p.Fields,
			},
		},
		Highlight: &types.Highlight{
			PreTags:  []string{"<mark>"},
			PostTags: []string{"</mark>"},
			Fields: map[string]types.HighlightField{
				"*": {},
			},
			RequireFieldMatch: &requireMatch,
		},
	}
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(req)
	return &buf
}
