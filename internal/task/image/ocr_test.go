package image

import (
	"encoding/base64"
	"os"
	"testing"
)

func TestProcessImageOCR(t *testing.T) {
	// 
}

func TestCallOCRAPI(t *testing.T) {
	// 读取测试图片文件
	imageFile := "../../../test.png" // 图片应放在项目根目录或适当的测试目录中
	
	// 打开并读取图片文件
	imgData, err := os.ReadFile(imageFile)
	if err != nil {
		t.Fatalf("无法读取测试图片文件: %v", err)
	}
	
	// 将图片转换为Base64编码
	base64Image := base64.StdEncoding.EncodeToString(imgData)
	
	// 调用OCR API
	resp, err := callOCRAPI(base64Image)
	if err != nil {
		t.Fatalf("OCR API调用失败: %v", err)
	}
	
	// 验证OCR API响应
	if resp.Errcode != 0 {
		t.Errorf("OCR API返回错误码: %d", resp.Errcode)
	}
	
	// 检查是否有识别结果
	if len(resp.OCRResponse) == 0 {
		t.Log("OCR未能识别任何文本")
	} else {
		t.Logf("成功识别出%d个文本区域", len(resp.OCRResponse))
		// 打印识别出的第一个文本，用于验证
		if len(resp.OCRResponse) > 0 {
			t.Logf("第一个识别文本: %s (置信度: %.2f%%)", 
				resp.OCRResponse[0].Text, resp.OCRResponse[0].Rate*100)
		}
	}
}