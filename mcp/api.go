package mcp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *MCPService) Chat(context string) (string, error) {
	// 提取信息
	requestBodyJSON, err := s.extract_keyword(context)
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return "", err
	}
	// 发送 POST 请求 读取响应体
	s.Context = context
	respBody, err := s.sendContextToModel(&requestBodyJSON)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return "", err
	}
	fmt.Println("响应内容:", string(*respBody))
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

func NewMCPService(name string, host string, key string) *MCPService {
	s := &MCPService{}
	s.Name = name
	s.Host = host
	s.Key = key
	s.Client = http.Client{}
	s.Client.Timeout = 60 * time.Second
	s.Client.Transport = &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       60 * time.Second,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
		ExpectContinueTimeout: 10 * time.Second,
	}
	return s
}

func (s *MCPService) AddTool(name string, description string, parameters Paramaters, handler func(args map[string]any) (string, error)) {
	tool := Tool{
		Type: "function",
		Function: struct {
			Name        string     `json:"name"`
			Description string     `json:"description"`
			Para        Paramaters `json:"parameters"`
		}{
			Name:        name,
			Description: description,
			Para:        parameters,
		},
		Handler: handler,
	}
	s.Tools = append(s.Tools, tool)
	fmt.Println("Tool", tool)
}
