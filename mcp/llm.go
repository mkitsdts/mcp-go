package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const request_content_type string = "application/json"

// 向大模型发送请求
func (s *MCPClient) send_request(data *[]byte) (*[]byte, error) {
	// 发送 POST 请求
	req, err := http.NewRequest("POST", s.host, bytes.NewBuffer(*data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", request_content_type)
	if s.key != "" {
		req.Header.Set("Authorization", "Bearer "+s.key)
	}
	fmt.Println("请求头:", req.Header)
	resp, err := s.GetHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return nil, err
	}
	return &respBody, nil
}

// 提取回答关键信息
func (s *MCPClient) create_request(context string, role string) (*[]byte, error) {
	if context != "" {
		s.context = append(s.context, req_mess{Role: role, Content: context})
	}
	contexts := make([]req_mess, len(s.context))
	copy(contexts, s.context)
	for key, value := range s.files {
		fmt.Println("文件名:", key)
		fmt.Println("文件内容:", value)
		contexts = append(contexts, req_mess{Role: "user", Content: "file<" + key + ">: " + value})
	}

	// 转换为JSON
	messagesJSON, err := json.Marshal(contexts)
	if err != nil {
		return nil, err
	}
	sum := len(s.tools)
	tools := make([]Tool, sum)
	// 将当前对话的工具添加到请求体
	copy(tools, s.tools)
	// 创建完整请求体
	requestBody := request{
		Model:       s.name,
		Messages:    json.RawMessage(messagesJSON),
		Temperature: 0.0,
		Stream:      false,
		Tools:       tools,
	}
	// 转换为JSON
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return nil, err
	}
	fmt.Println("Request body:", string(requestBodyJSON))
	return &requestBodyJSON, nil
}
