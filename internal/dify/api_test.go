package dify

import (
	"testing"
)

func TestCallWorkflow_Integration(t *testing.T) {

	// 创建测试请求
	req := &WorkflowRequest{
		Inputs: WorkflowInput{
			Action:  ActionExpand,
			Content: "斐波那契",
		},
	}

	// 调用 API
	resp, err := CallWorkflow(req)
	if err != nil {
		t.Fatalf("调用 API 失败: %v", err)
	}

	// 验证响应
	if resp.Status != "succeeded" {
		t.Errorf("期望状态为 succeeded，但得到: %v", resp.Status)
	}
	if resp.Message == "" {
		t.Error("响应消息为空")
	}

	// 打印响应以便调试
	t.Logf("API Status: %+v", resp.Status)
}
