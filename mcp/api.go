package mcp

import (
	"fmt"
	"net/http"
	"time"
)

func (s *MCPClient) Chat(context string) (string, error) {
	// 提取信息
	extractKeywordBodyJSON, err := s.create_request(context, "user")
	if err != nil {
		fmt.Println("转换请求体为JSON错误:", err)
		return "", err
	}
	// 发送 POST 请求
	keywordBody, err := s.send_request(extractKeywordBodyJSON)
	if err != nil {
		fmt.Println("读取响应体错误:", err)
		return "", err
	}
	fmt.Println("响应内容:", string(*keywordBody))
	// 解析结果并调用工具
	result, err := s.parseresp(keywordBody)
	if err != nil {
		if err.Error() == "error: no tool calls found in response" {
			fmt.Println("没有工具调用，直接返回响应内容")
			return result, nil
		}
		fmt.Println("解析响应结果错误:", err)
		return "", err
	}
	return result, nil
}

func NewMCPService(name string, host string, key string) *MCPService {
	s := &MCPService{}
	s.name = name
	s.host = host
	s.key = key
	s.Clients = map[string]MCPClient{}
	return s
}

func (s *MCPService) AddGlobalTool(name string, description string, parameters Paramaters, handler func(args map[string]any) (string, error)) {
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
	s.tools = append(s.tools, tool)
}

func (s *MCPClient) AddTool(name string, description string, parameters Paramaters, handler func(args map[string]any) (string, error)) {
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
	s.tools = append(s.tools, tool)
	fmt.Println("Tool", tool)
}

func (s *MCPService) NewClient(tag string) *MCPClient {
	c := MCPClient{}
	c.client = http.Client{}
	c.client.Timeout = 60 * time.Second
	c.client.Transport = &http.Transport{
		MaxIdleConns:          10,
		IdleConnTimeout:       60 * time.Second,
		DisableCompression:    true,
		ForceAttemptHTTP2:     true,
		ExpectContinueTimeout: 10 * time.Second,
	}
	c.context = []req_mess{}
	c.context = append(c.context, req_mess{Role: "system", Content: system_prompt})
	c.tools = []Tool{}
	c.host = &s.host
	c.key = &s.key
	c.name = &s.name
	c.golbaltool = &s.tools
	s.Clients[tag] = c
	return &c
}
