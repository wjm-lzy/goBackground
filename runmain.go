package main

import (
	"bytes"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"io"
	"net/http"
	"sync"
)

// 百度千帆配置
const (
	BAIDU_API_KEY = "bce-v3/ALTAK-xDqHxq05FYcbdbF8w06ZS/092b6de8d924a546fef661700b635a2b3f3387a6" // 替换为您的实际API Key
	CHAT_URL      = "https://qianfan.baidubce.com/v2/chat/completions"
	IMAGE_URL     = "https://qianfan.baidubce.com/v2/images" // 假设图像API使用相同格式
)

// 全局变量
var (
	imageTaskMap  sync.Map // 存储任务ID和图像URL
	taskExpiryMap sync.Map // 存储任务到期时间
	database      *sqlx.DB
)

// 百度千帆消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 百度千帆请求结构
type BaiduRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	System   string    `json:"system,omitempty"`
}

// 百度千帆响应结构
type BaiduResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	ErrorCode int    `json:"error_code,omitempty"`
	ErrorMsg  string `json:"error_msg,omitempty"`
}

// 调用百度千帆大模型
func askBaiduQianfan(systemPrompt string, messages []map[string]interface{}) (string, error) {
	// 转换为百度需要的消息格式
	baiduMessages := make([]Message, 0, len(messages))
	for _, msg := range messages {
		role, _ := msg["role"].(string)
		content, _ := msg["content"].(string)
		baiduMessages = append(baiduMessages, Message{
			Role:    role,
			Content: content,
		})
	}

	// 创建请求
	request := BaiduRequest{
		Model:    "ernie-3.5-8k", // 根据文档指定模型
		Messages: baiduMessages,
		System:   systemPrompt,
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %w", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", CHAT_URL, bytes.NewBuffer(payload))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置认证头 - 直接使用API Key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+BAIDU_API_KEY)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 处理响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 调试输出
	fmt.Printf("百度API响应: %s\n", string(body))

	var response BaiduResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 错误处理
	if response.ErrorCode != 0 {
		return "", fmt.Errorf("百度错误 %d: %s", response.ErrorCode, response.ErrorMsg)
	}

	return response.Choices[0].Message.Content, nil
}

// 测试大模型功能
//func testLLM() {
//	messages := []map[string]interface{}{
//		{"role": "user", "content": "你好，1到100的和是多少"},
//	}
//	response, err := askBaiduQianfan("你是一个有帮助的助手", messages)
//	if err != nil {
//		fmt.Printf("API错误: %v\n", err)
//	}
//	fmt.Printf("大模型测试响应: %s", response)
//}
