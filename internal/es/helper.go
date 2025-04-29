package es

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
	"unicode"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type SearchParams struct {
	Keyword string
	Fields  []string
	Index   string
}

func containsChinese(s string) bool {
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}

func BuildBody(p SearchParams) *bytes.Buffer {
	requireMatch := false

	// 判断 keyword 是否为中文
	var queryType *textquerytype.TextQueryType
	if containsChinese(p.Keyword) {
		t := textquerytype.Phrase
		queryType = &t
	} else {
		queryType = nil // 英文不设置 type，使用默认
	}

	req := search.Request{
		Query: &types.Query{
			MultiMatch: &types.MultiMatchQuery{
				Query:  p.Keyword,
				Fields: p.Fields,
				Type:   queryType,
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
