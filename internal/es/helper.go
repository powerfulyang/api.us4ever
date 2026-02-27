package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"unicode"

	"api.us4ever/internal/config"
	"api.us4ever/internal/logger"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

var (
	helperLogger *logger.Logger
)

func init() {
	var err error
	helperLogger, err = logger.New("helper")
	if err != nil {
		panic("failed to initialize helper logger: " + err.Error())
	}
}

const (
	vectorDims = 1024
)

type SearchParams struct {
	Keyword string
	Fields  []string
	Index   string
}

type EmbeddingReq struct {
	Text string `json:"text"`
}
type EmbeddingResp struct {
	Vector []float32 `json:"embedding"`
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

func Embed(ctx context.Context, text string) ([]float32, error) {
	// 超时 3 秒
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	appConfig := config.GetAppConfig()
	embedServiceURL := appConfig.Embedding.Endpoint

	reqBody, _ := json.Marshal(EmbeddingReq{Text: text})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, embedServiceURL, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		helperLogger.Errorw("embed service error", "err", err)
		// 返回一个模长极小的非零向量
		dummy := make([]float32, 1024)
		dummy[0] = 0.1
		return dummy, nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embed service %d: %s", resp.StatusCode, body)
	}
	var er EmbeddingResp
	if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
		return nil, err
	}
	// 简单校验维度
	if len(er.Vector) != vectorDims {
		return nil, fmt.Errorf("expect %d dims, got %d", vectorDims, len(er.Vector))
	}
	return er.Vector, nil
}

// MergeTextFields -------- 把多个字段拼出同一套 mapping --------
func MergeTextFields(names []string) map[string]any {
	out := make(map[string]any)
	for _, f := range names {
		out[f] = map[string]any{
			"type":            "text",
			"analyzer":        "ik_cjk",
			"search_analyzer": "ik_cjk",
			"fields": map[string]any{
				"ngram": map[string]any{
					"type":            "text",
					"analyzer":        "cjk_ngram_analyzer",
					"search_analyzer": "ik_cjk",
				},
			},
		}
	}
	return out
}
