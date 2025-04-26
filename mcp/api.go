package mcp

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (s *MCPService) Chat(context string) (string, error) {
	s.Messages = append(s.Messages, map[string]string{"role": "user", "content": context})
	// 提取信息
	requestBodyJSON, err := s.extract_keyword(context)
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return "", err
	}
	// 发送 POST 请求
	resp, err := s.Client.Post(s.Host+"/v1/chat/completions", "application/json", bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return "", err
	}
	fmt.Println("响应内容:", string(respBody))
	s.Messages = append(s.Messages, map[string]string{"role": "assistant", "content": string(respBody)})
	// 获取结果
	answer, err := s.get_answer(respBody)
	if err != nil {
		fmt.Println("解析响应结果错误:", err)
		return "", err
	}
	result, err := s.extract_result(answer)
	// 打印结果
	fmt.Println("结果:", result)
	if err != nil {
		fmt.Println("解析结果错误:", err)
		return "", err
	}
	return result, nil
}

func NewMCPService(name string, host string) *MCPService {
	s := &MCPService{}
	s.Name = name
	s.Host = host
	s.Client = http.Client{}
	s.Client.Timeout = 60 * time.Second
	s.Client.Transport = &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       60 * time.Second,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
		ExpectContinueTimeout: 10 * time.Second,
	}
	if !s.try_get_model_info() {
		panic("Failed to connect to model server")
	} else {
		fmt.Println("Model server connected successfully")
	}
	s.Messages = append(s.Messages, map[string]string{"role": "system", "content": system_prompt})
	return s
}

func (s *MCPService) AddTool(name string, description string, parameters Para, handler func(args map[string]any) (string, error)) {
	tool := Tool{
		Type: "function",
		Function: struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			Parameters  Para   `json:"parameters"`
		}{
			Name:        name,
			Description: description,
			Parameters:  parameters,
		},
		Handler: handler,
	}
	s.Tools = append(s.Tools, tool)
	fmt.Println("Tool", tool)
}
