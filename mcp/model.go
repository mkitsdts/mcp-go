package mcp

import "net/http"

type Para struct {
	Type       string         `json:"type"`
	Properties map[string]any `json:"properties"`
	Required   []string       `json:"required"`
}

type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Parameters  Para   `json:"parameters"`
	} `json:"function"`
	Handler func(args map[string]any) (string, error) `json:"-"`
}

// 核心服务
// 包含上下文协议、工具、消息等
type MCPService struct {
	Client  http.Client
	Host    string
	Name    string
	Key     string
	Tools   []Tool
	Context string
}
