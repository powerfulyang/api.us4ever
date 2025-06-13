package dify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"bufio"

	"api.us4ever/internal/config"
	"api.us4ever/internal/logger"
	"go.uber.org/zap"
)

var (
	difyLogger *logger.Logger
)

func init() {
	var err error
	difyLogger, err = logger.New("dify")
	if err != nil {
		panic("failed to initialize dify logger: " + err.Error())
	}
}

type ResponseMode string

const (
	ResponseModeStreaming ResponseMode = "streaming"
	ResponseModeBlocking  ResponseMode = "blocking"
)

// ActionType 定义可用的 action 类型
type ActionType string

const (
	ActionTitle   ActionType = "title"
	ActionContent ActionType = "content"
	ActionExpand  ActionType = "expand"
)

// WorkflowInput 定义 inputs 字段的结构
type WorkflowInput struct {
	// Action 指定要执行的操作类型：title、content 或 expand
	Action ActionType `json:"action"`

	// Content 是操作对应的文本内容
	Content string `json:"content"`
}

type WorkflowRequest struct {
	Inputs       WorkflowInput `json:"inputs"`
	ResponseMode ResponseMode  `json:"response_mode"`
	User         string        `json:"user"`
}

func (r *WorkflowRequest) SetDefaults() {
	if r.ResponseMode == "" {
		r.ResponseMode = ResponseModeStreaming
	}
	if r.User == "" {
		r.User = "default"
	}
}

// WorkflowResponse 定义了API响应的结构
type WorkflowResponse struct {
	TaskID        string `json:"task_id"`
	WorkflowRunID string `json:"workflow_run_id"`
	Data          struct {
		ID         string `json:"id"`
		WorkflowID string `json:"workflow_id"`
		Status     string `json:"status"`
		Outputs    struct {
			Text string `json:"text"`
		} `json:"outputs"`
		Error       interface{} `json:"error"`
		ElapsedTime float64     `json:"elapsed_time"`
		TotalTokens int         `json:"total_tokens"`
		TotalSteps  int         `json:"total_steps"`
		CreatedAt   int64       `json:"created_at"`
		FinishedAt  int64       `json:"finished_at"`
	} `json:"data"`
}

// WorkflowStreamResponse 定义了流式响应的结构
type WorkflowStreamResponse struct {
	Event         string `json:"event"`
	WorkflowRunID string `json:"workflow_run_id"`
	TaskID        string `json:"task_id"`
	Data          struct {
		ID         string `json:"id"`
		WorkflowID string `json:"workflow_id"`
		Status     string `json:"status"`
		Outputs    struct {
			Text string `json:"text"`
		} `json:"outputs"`
		Error                interface{} `json:"error"`
		ElapsedTime          float64     `json:"elapsed_time"`
		TotalTokens          int         `json:"total_tokens"`
		TotalSteps           int         `json:"total_steps"`
		CreatedAt            int64       `json:"created_at"`
		FinishedAt           int64       `json:"finished_at"`
		Text                 string      `json:"text"`
		FromVariableSelector []string    `json:"from_variable_selector,omitempty"`
	} `json:"data"`
}

// WorkflowResult 定义了统一的返回结果结构
type WorkflowResult struct {
	Message string
	Status  string
	Error   error
}

// CallWorkflow 调用 Dify Workflow API
func CallWorkflow(req *WorkflowRequest) (*WorkflowResult, error) {
	req.SetDefaults()
	// 从配置中获取 endpoint 和 apiKey
	appConfig := config.GetAppConfig()
	if appConfig == nil {
		return nil, fmt.Errorf("无法获取应用配置")
	}

	if appConfig.Dify.Endpoint == "" || appConfig.Dify.ApiKey == "" {
		return nil, fmt.Errorf("dify 配置不完整: endpoint 或 apiKey为空")
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求数据失败: %v", err)
	}

	httpReq, err := http.NewRequest("POST", appConfig.Dify.Endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+appConfig.Dify.ApiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	// httpReq.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用 dify 失败: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			difyLogger.Warn("failed to close response body",
				zap.Error(err),
			)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
	}

	result := &WorkflowResult{}

	switch req.ResponseMode {
	case ResponseModeStreaming:
		// 处理 SSE 流式响应
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, fmt.Errorf("读取SSE响应失败: %v", err)
			}

			line = strings.TrimSpace(line)

			// 跳过空行
			if line == "" {
				continue
			}

			// 解析 SSE 数据行
			if strings.HasPrefix(line, "data: ") {
				jsonPart := strings.TrimPrefix(line, "data: ")
				var streamResp WorkflowStreamResponse
				if err := json.Unmarshal([]byte(jsonPart), &streamResp); err != nil {
					return nil, fmt.Errorf("解析SSE数据失败: %v", err)
				}

				// 只处理 text_chunk 事件
				if streamResp.Event == "text_chunk" {
					result.Message += streamResp.Data.Text
				}

				if streamResp.Event == "workflow_finished" {
					result.Status = streamResp.Data.Status
				}
			}
		}

	case ResponseModeBlocking:
		// 处理阻塞式响应
		var blockResp WorkflowResponse
		if err := json.NewDecoder(resp.Body).Decode(&blockResp); err != nil {
			return nil, fmt.Errorf("解析API响应失败: %v", err)
		}
		result.Message = blockResp.Data.Outputs.Text
		result.Status = blockResp.Data.Status

	default:
		return nil, fmt.Errorf("不支持的响应模式: %s", req.ResponseMode)
	}

	return result, nil
}
